package proxy

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/lflxp/lflxp-k8s/utils"

	"github.com/gin-gonic/gin"
)

var rancheUnauthedApis = []string{"v3/connect", "v3/settings", "ping", "cacerts", "system-agent-install.sh", "healthz", "assets"}

func setTokenToUrl(target string, rawUrl *url.URL) error {
	// 这边是从设置里拿代理值
	//equipment, err := proxy.GetEquipment()
	//if err != nil {
	//	return err
	//}

	proxyUrl, err := url.Parse(target)
	if err != nil {
		return err
	}

	tpaths := strings.Split(target, "/")
	rawpaths := strings.Split(rawUrl.Path, "/")

	result := []string{}
	if len(tpaths) == 4 {
		result = []string{"", tpaths[len(tpaths)-1]}
		result = append(result, rawpaths[2:]...)
	} else if len(tpaths) > 4 {
		t := len(tpaths) - 3
		result = []string{"", tpaths[len(tpaths)-t]}
		result = append(result, rawpaths[t:]...)
	}

	rawUrl.Path = strings.Join(result, "/")
	rawUrl.Scheme = proxyUrl.Scheme
	rawUrl.Host = proxyUrl.Host
	// ruq := rawUrl.Query()
	// ruq.Add("token", token)
	// rawUrl.RawQuery = ruq.Encode()
	return nil
}

// 自定义http七层代理服务
func NewHttpProxyByGinCustomRaw(target string, filter map[string]string) func(c *gin.Context) {
	return func(c *gin.Context) {
		if len(filter) > 0 {
			for key, value := range filter {
				if c.GetHeader(key) != value {
					utils.SendErrorMessage(c, http.StatusBadRequest, "GetHeader", fmt.Sprintf("Header %s:%s not define", key, value))
					return
				}
			}
		}

		err := setTokenToUrl(target, c.Request.URL)
		if err != nil {
			utils.SendErrorMessage(c, http.StatusInternalServerError, "setTokenToUrl", fmt.Sprintf("填写的地址有误: %s", err.Error()))
			c.Abort()
			return
		}

		req, err := http.NewRequestWithContext(c, c.Request.Method, c.Request.URL.String(), c.Request.Body)
		if err != nil {
			utils.SendErrorMessage(c, http.StatusInternalServerError, "NewRequestWithContext", err.Error())
			c.Abort()
			return
		}
		defer req.Body.Close()
		req.Header = c.Request.Header
		if c.Request.Header.Get("token") != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Request.Header.Get("token")))
		}

		if len(c.Request.Cookies()) > 0 {
			for _, cookie := range c.Request.Cookies() {
				req.AddCookie(cookie)
			}
		}

		if clientIP, _, err := net.SplitHostPort(c.Request.RemoteAddr); err == nil {
			prior, ok := req.Header["X-Forwarded-For"]
			omit := ok && prior == nil // Issue 38079: nil now means don't populate the header
			if len(prior) > 0 {
				clientIP = strings.Join(prior, ", ") + ", " + clientIP
			}
			if !omit {
				req.Header.Set("X-Forwarded-For", clientIP)
			}
		}

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			utils.SendErrorMessage(c, http.StatusInternalServerError, "DefaultClient.Do()", err.Error())
			c.Abort()
			return
		}

		extraHeaders := make(map[string]string)
		extraHeaders["PROXY"] = "caas-proxy"

		// 放开第三方UI代理
		if strings.Contains(c.Request.URL.Path, "monitoring-system") {
			for k := range resp.Header {
				for j := range resp.Header[k] {
					c.Header(k, resp.Header[k][j])
				}
			}

			for _, cookie := range resp.Cookies() {
				c.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
			}

			if strings.Contains(c.Request.URL.String(), "proxy/") &&
				!strings.Contains(c.Request.URL.String(), "/clusters/local/") &&
				strings.Contains(resp.Header.Get("Content-Type"), "text/html") {

				needReplaceStrHttp := "http://" + c.Request.URL.Host
				needReplaceStrHttps := "https://" + c.Request.URL.Host
				slog.Info("needReplaceStrHttp: %s", needReplaceStrHttp)
				slog.Info("needReplaceStrHttps: %s", needReplaceStrHttps)

				if resp.Header.Get("Content-Encoding") == "gzip" {
					orgBodyGzipReader, err := gzip.NewReader(resp.Body)
					if err != nil {
						utils.SendErrorMessage(c, http.StatusInternalServerError, "gzip.NewReader(resp.Body)", err.Error())
						c.Abort()
						return
					}

					orgBodyBuf := new(bytes.Buffer)
					_, err = orgBodyBuf.ReadFrom(orgBodyGzipReader)
					if err != nil {
						utils.SendErrorMessage(c, http.StatusInternalServerError, "orgBodyBuf.ReadFrom(orgBodyGzipReader)", err.Error())
						c.Abort()
						return
					}

					orgBodyBytes := orgBodyBuf.Bytes()

					err = resp.Body.Close()
					if err != nil {
						utils.SendErrorMessage(c, http.StatusInternalServerError, "resp.Body.Close()", err.Error())
						c.Abort()
						return
					}

					orgBodyBytes = bytes.Replace(orgBodyBytes, []byte(needReplaceStrHttp), []byte(""), -1)  // replace html
					orgBodyBytes = bytes.Replace(orgBodyBytes, []byte(needReplaceStrHttps), []byte(""), -1) // replace html

					var newBodyBuf bytes.Buffer
					newBodyWriter := gzip.NewWriter(&newBodyBuf)
					defer newBodyWriter.Close()

					_, err = newBodyWriter.Write(orgBodyBytes)
					if err != nil {
						utils.SendErrorMessage(c, http.StatusInternalServerError, "newBodyWriter.Write(orgBodyBytes)", err.Error())
						c.Abort()
						return
					}

					err = newBodyWriter.Flush()
					if err != nil {
						utils.SendErrorMessage(c, http.StatusInternalServerError, "newBodyWriter.Flush()", err.Error())
						c.Abort()
						return
					}

					newBodyBytes := newBodyBuf.Bytes()
					c.DataFromReader(resp.StatusCode, int64(len(newBodyBytes)), resp.Header.Get("Content-AlertType"), bytes.NewReader(newBodyBytes), extraHeaders)
				} else {

					b, _ := ioutil.ReadAll(resp.Body)
					b = bytes.Replace(b, []byte(needReplaceStrHttp), []byte(""), -1)  // replace html
					b = bytes.Replace(b, []byte(needReplaceStrHttps), []byte(""), -1) // replace html

					c.DataFromReader(resp.StatusCode, int64(len(b)), resp.Header.Get("Content-AlertType"), bytes.NewReader(b), extraHeaders)
				}

			} else {
				c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-AlertType"), resp.Body, extraHeaders)
			}

		} else {
			var msg interface{}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				slog.Error(err.Error())
			}
			err = json.Unmarshal(body, &msg)
			if err != nil {
				if strings.Contains(string(body), "not authorized") {
					utils.SendErrorMessage(c, 401, "Unauthorized", string(body))
				} else if resp.StatusCode == 422 {
					info, err := utils.ParseE422(body)
					if err != nil {
						utils.SendErrorMessage(c, resp.StatusCode, string(resp.StatusCode), err.Error())
					} else {
						utils.SendErrorMessage(c, resp.StatusCode, string(resp.StatusCode), info.String())
					}
				} else {
					utils.SendErrorMessage(c, resp.StatusCode, "json.Unmarshal", fmt.Sprintf("%d body %s error: %s", resp.StatusCode, string(body), err.Error()))
				}

				c.Abort()
				return
			}

			if resp.StatusCode < 400 {
				// header 也带过来
				for k := range resp.Header {
					for j := range resp.Header[k] {
						c.Header(k, resp.Header[k][j])
					}
				}

				for _, cookie := range resp.Cookies() {
					c.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, c.Request.Host, cookie.Secure, cookie.HttpOnly)
				}

				// c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-AlertType"), resp.Body, extraHeaders)
				utils.SendSuccessMessage(c, resp.StatusCode, msg)
				c.Abort()
			} else {
				utils.SendErrorMessage(c, resp.StatusCode, "error", string(body))
				c.Abort()
			}
		}
	}
}

func needHandleGrafanaHtmlResponse(c *gin.Context, resp *http.Response) bool {
	// return strings.Contains(c.Request.URL.Path, "monitoring-system") &&
	// 	strings.Contains(c.Request.URL.String(), "proxy/") &&
	// 	!strings.Contains(c.Request.URL.String(), "/clusters/local/") &&
	// 	strings.Contains(resp.Header.Get("Content-Type"), "text/html")
	return strings.Contains(c.Request.URL.Path, "/grafana") &&
		strings.Contains(resp.Header.Get("Content-Type"), "text/html")
}

func IsRancheUnauthedApis(c *gin.Context) bool {
	for _, url := range rancheUnauthedApis {
		if strings.Contains(c.Request.URL.String(), url) {
			return true
		}
	}
	return false
}

// 处理grafana ui query接口
func handleGrafanaQueryRequest(c *gin.Context, req *http.Request) {
	if !strings.Contains(c.Request.URL.Path, "/grafana/api") {
		return
	}

	if c.Request.Header.Get("Accept") != "" {
		req.Header.Set("Accept", c.Request.Header.Get("Accept"))
	}

	if c.Request.Header.Get("Accept-Encoding") != "" {
		req.Header.Set("Accept-Encoding", c.Request.Header.Get("Accept-Encoding"))
	}

	if c.Request.Header.Get("Accept-Language") != "" {
		req.Header.Set("Accept-Language", c.Request.Header.Get("Accept-Language"))
	}

	if c.Request.Header.Get("Content-Length") != "" {
		req.Header.Set("Content-Length", c.Request.Header.Get("Content-Length"))
	}

	if c.Request.Header.Get("Content-Type") != "" {
		req.Header.Set("Content-Type", c.Request.Header.Get("Content-Type"))
	}
}

func handleGrafanaHtmlResponse(c *gin.Context, resp *http.Response) {
	extraHeaders := make(map[string]string)
	extraHeaders["PROXY"] = "caas-proxy"

	needReplaceStrHttp := "http://" + c.Request.URL.Host
	needReplaceStrHttps := "https://" + c.Request.URL.Host
	slog.Info("needReplaceStrHttp: %s", needReplaceStrHttp)
	slog.Info("needReplaceStrHttps: %s", needReplaceStrHttps)

	if resp.Header.Get("Content-Encoding") == "gzip" {
		orgBodyGzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			utils.SendErrorMessage(c, http.StatusInternalServerError, "gzip.NewReader(resp.Body)", err.Error())
			c.Abort()
			return
		}

		orgBodyBuf := new(bytes.Buffer)
		_, err = orgBodyBuf.ReadFrom(orgBodyGzipReader)
		if err != nil {
			utils.SendErrorMessage(c, http.StatusInternalServerError, "orgBodyBuf.ReadFrom(orgBodyGzipReader)", err.Error())
			c.Abort()
			return
		}

		orgBodyBytes := orgBodyBuf.Bytes()

		err = resp.Body.Close()
		if err != nil {
			utils.SendErrorMessage(c, http.StatusInternalServerError, "resp.Body.Close()", err.Error())
			c.Abort()
			return
		}

		orgBodyBytes = bytes.Replace(orgBodyBytes, []byte(needReplaceStrHttp), []byte(""), -1)  // replace html
		orgBodyBytes = bytes.Replace(orgBodyBytes, []byte(needReplaceStrHttps), []byte(""), -1) // replace html

		var newBodyBuf bytes.Buffer
		newBodyWriter := gzip.NewWriter(&newBodyBuf)
		defer newBodyWriter.Close()

		_, err = newBodyWriter.Write(orgBodyBytes)
		if err != nil {
			utils.SendErrorMessage(c, http.StatusInternalServerError, "newBodyWriter.Write(orgBodyBytes)", err.Error())
			c.Abort()
			return
		}

		err = newBodyWriter.Flush()
		if err != nil {
			utils.SendErrorMessage(c, http.StatusInternalServerError, "newBodyWriter.Flush()", err.Error())
			c.Abort()
			return
		}

		newBodyBytes := newBodyBuf.Bytes()
		c.DataFromReader(resp.StatusCode, int64(len(newBodyBytes)), resp.Header.Get("Content-Type"), bytes.NewReader(newBodyBytes), extraHeaders)
	} else {

		b, _ := ioutil.ReadAll(resp.Body)
		b = bytes.Replace(b, []byte(needReplaceStrHttp), []byte(""), -1)  // replace html
		b = bytes.Replace(b, []byte(needReplaceStrHttps), []byte(""), -1) // replace html

		c.DataFromReader(resp.StatusCode, int64(len(b)), resp.Header.Get("Content-AlertType"), bytes.NewReader(b), extraHeaders)
	}
}

// 自定义http七层代理服务
// 原生代理 不做任何修数
func NewHttpProxyByGinCustom(target string, filter map[string]string) func(c *gin.Context) {
	return func(c *gin.Context) {
		if len(filter) > 0 {
			for key, value := range filter {
				if c.GetHeader(key) != value {
					utils.SendErrorMessage(c, http.StatusBadRequest, "GetHeader", fmt.Sprintf("Header %s:%s not define", key, value))
					return
				}
			}
		}

		// // check the proxy request whether it is websocket
		if c.Request.Header.Get("Connection") == "Upgrade" && c.Request.Header.Get("Upgrade") == "websocket" {
			WebSocketProxy2(target, c)
			return
		}

		// if IsRancheUnauthedApis(c) {
		// 	RancherApiProxy(target, c)
		// 	return
		// }

		err := setTokenToUrl(target, c.Request.URL)
		if err != nil {
			utils.SendErrorMessage(c, http.StatusInternalServerError, "setTokenToUrl", fmt.Sprintf("填写的地址有误: %s", err.Error()))
			c.Abort()
			return
		}

		req, err := http.NewRequestWithContext(c, c.Request.Method, c.Request.URL.String(), c.Request.Body)
		if err != nil {
			utils.SendErrorMessage(c, http.StatusInternalServerError, "NewRequestWithContext", err.Error())
			c.Abort()
			return
		}
		defer req.Body.Close()
		// BUG: 前端header透传回导致接口权限报错
		// req.Header = c.Request.Header
		// if c.Request.Header.Get("token") != "" {
		// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Request.Header.Get("token")))
		// }

		handleGrafanaQueryRequest(c, req)

		if len(c.Request.Cookies()) > 0 {
			for _, cookie := range c.Request.Cookies() {
				req.AddCookie(cookie)
			}
		}

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			utils.SendErrorMessage(c, http.StatusInternalServerError, "DefaultClient.Do()", err.Error())
			c.Abort()
			return
		}
		defer resp.Body.Close()

		extraHeaders := make(map[string]string)
		extraHeaders["PROXY"] = "caas-proxy"

		// header 也带过来
		for k := range resp.Header {
			for j := range resp.Header[k] {
				c.Header(k, resp.Header[k][j])
			}
		}

		for _, cookie := range resp.Cookies() {
			c.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, c.Request.Host, cookie.Secure, cookie.HttpOnly)
		}

		if needHandleGrafanaHtmlResponse(c, resp) {
			handleGrafanaHtmlResponse(c, resp)
			return
		}

		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, extraHeaders)
		// utils.SendSuccessMessage(c, resp.StatusCode, msg)
		c.Abort()
	}
}

func WebSocketProxy2(target string, c *gin.Context) {
	// parsedUrl, err := url.Parse(target)
	// if err != nil {
	// 	logrus.Errorf("error to parse the url in the proxyRequest method, proxy path: %s, error: %v", target, err)
	// 	return
	// }

	// var tlsConfig = &tls.Config{
	// 	InsecureSkipVerify: true,
	// }

	// var transport http.RoundTripper = &http.Transport{
	// 	Proxy: nil,
	// 	DialContext: (&net.Dialer{
	// 		Timeout:   30 * time.Second,
	// 		KeepAlive: 30 * time.Second,
	// 	}).DialContext,
	// 	MaxIdleConns:          100,
	// 	IdleConnTimeout:       90 * time.Second,
	// 	TLSHandshakeTimeout:   10 * time.Second,
	// 	ExpectContinueTimeout: 1 * time.Second,
	// 	TLSClientConfig:       tlsConfig,
	// 	DisableCompression:    true,
	// }

	// parsedUrl.Path = c.Request.URL.Path
	// proxy := httputil.NewSingleHostReverseProxy(parsedUrl)
	// proxy.Transport = transport
	// proxy.ModifyResponse = func(r *http.Response) error {
	// 	r.Header.Del("X-Frame-Options")
	// 	return nil
	// }

	// c.Request.Header.Del("User-Agent")
	// c.Request.Header.Del("Origin")

	// err = setTokenToUrl(target, c.Request.URL)
	// if err != nil {
	// 	log.Error(err)
	// 	return
	// }
	// proxy.ServeHTTP(c.Writer, c.Request)
	c.JSONP(200, gin.H{"hello": "world"})
}

// WebSocketProxy use httputil.ReverseProxy proxy websocket
func WebSocketProxy(target string, c *gin.Context) {
	parsedUrl, err := url.Parse(target)
	if err != nil {
		slog.Error("error to parse the url in the proxyRequest method, proxy", "PATH", target, "ERROR", err.Error())
		return
	}
	director := func(req *http.Request) {
		req.URL.Scheme = parsedUrl.Scheme
		req.URL.Host = parsedUrl.Host
	}

	proxy := &httputil.ReverseProxy{
		FlushInterval: 500 * time.Millisecond,
		Director:      director,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	c.Request.Header.Del("User-Agent")
	c.Request.Header.Del("Origin")

	// cookie, err2 := middlewares.ParseToken(c)
	// if err2 != nil {
	// 	logrus.Errorf("error to parse the cookie error:%s", err.Error())
	// 	return
	// }
	// c.Request.Header.Set("token", cookie.Sess)
	// // TODO: rancherv2.6.2不支持未赋予work角色单节点集群执行kubectl命令行，临时手动设置kubect-pod容忍度待rancher官方bug修复
	// if clusterID := checkKubectlShell(c); clusterID != "" {
	// 	stop := make(chan bool)
	// 	go modifyKubectlPod(c, clusterID, stop)
	// 	proxy.ServeHTTP(c.Writer, c.Request)
	// 	stop <- true
	// } else {
	// 	proxy.ServeHTTP(c.Writer, c.Request)
	// }

	proxy.ServeHTTP(c.Writer, c.Request)
}

// // NewHttpReverseProxy use httputil.ReverseProxy realize the proxy
// func NewHttpReverseProxy(target string, filter map[string]string) func(c *gin.Context) {
// 	return func(c *gin.Context) {
// 		if len(filter) > 0 {
// 			for key, value := range filter {
// 				if c.GetHeader(key) != value {
// 					utils.SendErrorMessage(c, http.StatusBadRequest, "GetHeader", fmt.Sprintf("Header %s:%s not define", key, value))
// 					return
// 				}
// 			}
// 		}

// 		// TODO:修改透传代理，解决代理rancher接口返回httpCode != 200时,panic: http: wrote more than the declared Content-Length的问题
// 		parsedUrl, err := url.Parse(target)
// 		if err != nil {
// 			logrus.Errorf("error to parse the url in the proxyRequest method, proxy path: %s, error: %v", target, err)
// 			return
// 		}
// 		director := func(req *http.Request) {
// 			req.URL.Scheme = parsedUrl.Scheme
// 			req.URL.Host = parsedUrl.Host
// 		}

// 		proxy := &httputil.ReverseProxy{
// 			FlushInterval: 500 * time.Millisecond,
// 			Director:      director,
// 			Transport: &http.Transport{
// 				TLSClientConfig: &tls.Config{
// 					InsecureSkipVerify: true,
// 				},
// 			},
// 		}

// 		c.Request.Header.Del("User-Agent")
// 		c.Request.Header.Del("Origin")

// 		cookie, err2 := middlewares.ParseToken(c)
// 		if err2 != nil {
// 			logrus.Errorf("error to parse the cookie error:%s", err.Error())
// 			return
// 		}
// 		c.Request.Header.Set("token", cookie.Sess)
// 		// TODO: rancherv2.6.2不支持未赋予work角色单节点集群执行kubectl命令行，临时手动设置kubect-pod容忍度待rancher官方bug修复
// 		if clusterID := checkKubectlShell(c); clusterID != "" {
// 			stop := make(chan bool)
// 			go modifyKubectlPod(c, clusterID, stop)
// 			proxy.ServeHTTP(c.Writer, c.Request)
// 			stop <- true
// 		} else {
// 			proxy.ServeHTTP(c.Writer, c.Request)
// 		}

// 		// check the proxy request whether need to return to the unified structure
// 		if c.Request.Header.Get("Connection") == "Upgrade" && c.Request.Header.Get("Upgrade") == "websocket" {
// 			c.Request.Header.Del("Origin")

// 			cookie, err2 := middlewares.ParseToken(c)
// 			if err2 != nil {
// 				logrus.Errorf("error to parse the cookie error:%s", err.Error())
// 				return
// 			}
// 			c.Request.Header.Set("token", cookie.Sess)
// 			// TODO: rancherv2.6.2不支持未赋予work角色单节点集群执行kubectl命令行，临时手动设置kubect-pod容忍度待rancher官方bug修复
// 			if clusterID := checkKubectlShell(c); clusterID != "" {
// 				stop := make(chan bool)
// 				go modifyKubectlPod(c, clusterID, stop)
// 				proxy.ServeHTTP(c.Writer, c.Request)
// 				stop <- true
// 			} else {
// 				proxy.ServeHTTP(c.Writer, c.Request)
// 			}
// 		} else {
// 			// proxy only support Accept=application/json
// 			c.Request.Header.Set("Accept", "application/json")
// 			rpw := NewRespProxyWriter(c, true)
// 			proxy.ServeHTTP(rpw, c.Request)
// 			rpw.FlushBody()
// 		}
// 	}
// }

// // RancherApiProxy 用于无鉴权的rancher api通过apis代理
// func RancherApiProxy(target string, ctx *gin.Context) {
// 	parsedUrl, err := url.Parse(target)
// 	if err != nil {
// 		logrus.Errorf("error to parse the url in the proxyRequest method, proxy path: %s, error: %v", target, err)
// 		return
// 	}

// 	director := func(req *http.Request) {
// 		req.URL.Scheme = parsedUrl.Scheme
// 		req.URL.Host = parsedUrl.Host
// 		req.URL.Path = strings.TrimPrefix(ctx.Request.URL.Path, unAuthUrl)
// 	}

// 	proxy := &httputil.ReverseProxy{
// 		Director: director,
// 		Transport: &http.Transport{
// 			TLSClientConfig: &tls.Config{
// 				InsecureSkipVerify: true,
// 			},
// 		},
// 	}
// 	proxy.ServeHTTP(ctx.Writer, ctx.Request)
// }

// // ResponseProxyWriter extends the http.ResponseWriter,
// type ResponseProxyWriter interface {
// 	http.ResponseWriter
// 	FlushBody()
// }

// // ResponseProxyWriter extends the http.ResponseWriter,
// type responseProxyWriter struct {
// 	context        *gin.Context
// 	output         []byte
// 	needModifyBody bool
// }

// func NewRespProxyWriter(c *gin.Context, needModifyBody bool) ResponseProxyWriter {
// 	return &responseProxyWriter{
// 		context:        c,
// 		needModifyBody: needModifyBody,
// 	}
// }

// func (rpw *responseProxyWriter) Header() http.Header {
// 	return rpw.context.Writer.Header()
// }

// func (rpw *responseProxyWriter) Write(bytes []byte) (int, error) {
// 	if !rpw.needModifyBody {
// 		return rpw.context.Writer.Write(bytes)
// 	}

// 	rpw.output = append(rpw.output, bytes...)
// 	return len(bytes), nil
// }

// func (rpw *responseProxyWriter) WriteHeader(statusCode int) {
// 	rpw.context.Writer.WriteHeader(statusCode)
// }

// func (rpw *responseProxyWriter) FlushBody() {
// 	if !rpw.needModifyBody {
// 		rpw.context.Writer.Flush()
// 		return
// 	}

// 	// 1、判断是否是Content-Encoding==gzip是则进行解压处理
// 	if rpw.Header().Get("Content-Encoding") == "gzip" {
// 		respBody, err := gzip.NewReader(bytes.NewReader(rpw.output))
// 		if err != nil {
// 			logrus.Errorf("unpack response body error: %s", err.Error())
// 			utils.SendErrorMessage(rpw.context, http.StatusBadRequest, utils.SystemError, "")
// 			return
// 		}

// 		body, err := ioutil.ReadAll(respBody)
// 		if err != nil {
// 			utils.SendErrorMessage(rpw.context, http.StatusBadRequest, utils.SystemError, "")
// 			return
// 		}
// 		rpw.output = body
// 	}

// 	// 2、将解压的body封装到统一返回结构体中
// 	var msg interface{}
// 	err := json.Unmarshal(rpw.output, &msg)
// 	if err != nil {
// 		logrus.Errorf("json Unmarshal error: %s", err.Error())
// 		utils.SendErrorMessage(rpw.context, http.StatusBadRequest, utils.SystemError, "")
// 		return
// 	}

// 	info := &utils.Result{
// 		Success: true,
// 		Data:    msg,
// 		Host:    rpw.context.Request.URL.Path,
// 	}
// 	data, err := json.Marshal(info)
// 	if err != nil {
// 		logrus.Errorf("json Marshal error: %s", err.Error())
// 		utils.SendErrorMessage(rpw.context, http.StatusBadRequest, utils.SystemError, "")
// 		return
// 	}

// 	// 3、再将统一封装的数据进行压缩
// 	if rpw.Header().Get("Content-Encoding") == "gzip" {
// 		var zBuf bytes.Buffer
// 		zw := gzip.NewWriter(&zBuf)
// 		if _, err = zw.Write(data); err != nil {
// 			logrus.Errorf("pack response body error: %s", err.Error())
// 			utils.SendErrorMessage(rpw.context, http.StatusBadRequest, utils.SystemError, "")
// 			return
// 		}
// 		zw.Close()
// 		data = zBuf.Bytes()
// 	}

// 	rpw.context.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
// 	if _, err := rpw.context.Writer.Write(data); err == nil {
// 		rpw.context.Writer.Flush()
// 	} else {
// 		logrus.Errorf("write data to response error: %s", err.Error())
// 		utils.SendErrorMessage(rpw.context, http.StatusBadRequest, utils.SystemError, "")
// 		return
// 	}
// }

// // checkKubectlShell determine whether kubectl-shell request, is it return clusterID
// func checkKubectlShell(c *gin.Context) string {
// 	if strings.Contains(c.Request.URL.Path, "v3/clusters") && c.Query("shell") != "" {
// 		clusterID := strings.Split(c.Request.URL.Path, "/")[3]
// 		if clusterID == "local" {
// 			return ""
// 		}
// 		return clusterID
// 	}
// 	return ""
// }

// //  modifyKubectlPod modify dashboard-shell-pod toleration, support dashboard-shell-pod deploy to control-plane/etcd node
// func modifyKubectlPod(ctx *gin.Context, clusterID string, stop chan bool) {
// 	tolerationSeconds := int64(6000)
// 	ticker := time.NewTicker(5 * time.Second)

// 	for {
// 		select {
// 		case <-stop:
// 			return
// 		case <-ticker.C:
// 			pods, err := pod.ListNativePods(ctx, clusterID)
// 			if err != nil {
// 				logrus.Errorf("modify shell-pod tolerations error: %s", err.Error())
// 				return
// 			}

// 			for _, updatePod := range pods.Data {
// 				if strings.Contains(updatePod.Name, "dashboard-shell") {
// 					if updatePod.Labels["modified-toleration"] != "" {
// 						continue
// 					}

// 					if updatePod.Labels == nil {
// 						updatePod.Labels = map[string]string{}
// 					}
// 					updatePod.Labels["modified-toleration"] = "true"

// 					newTolerations := []v1.Toleration{
// 						{
// 							Key:      "node-role.kubernetes.io/control-plane",
// 							Operator: "Exists",
// 							Effect:   "NoSchedule",
// 						},
// 						{
// 							Key:               "node-role.kubernetes.io/etcd",
// 							Operator:          "Exists",
// 							Effect:            "NoExecute",
// 							TolerationSeconds: &tolerationSeconds,
// 						},
// 					}
// 					updatePod.Spec.Tolerations = append(updatePod.Spec.Tolerations, newTolerations...)
// 					err := pod.UpdatePod(ctx, clusterID, updatePod)
// 					if err != nil {
// 						logrus.Errorf("update cluster[%s] kubectl shell pod[%s] error: %s", clusterID, updatePod.Name, err.Error())
// 						return
// 					}
// 				}
// 			}
// 		}
// 	}
// }
