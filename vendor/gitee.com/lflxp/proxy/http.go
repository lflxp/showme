package proxy

import (
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

// http 七层代理
// 参数： https://www.baidu.com
func NewHttpProxyByGin(target string, filter map[string]string) func(c *gin.Context) {
	return func(c *gin.Context) {
		// 判断header
		if len(filter) > 0 {
			for key, value := range filter {
				if c.GetHeader(key) != value {
					c.String(http.StatusBadRequest, fmt.Sprintf("Header %s:%s not define", key, value))
					return
				}
			}
		}

		// var proxyUrl = new(url.URL)
		// proxyUrl.Scheme = "https"
		// proxyUrl.Host = "localhost"
		// proxyUrl.RawQuery = url.QueryEscape("proxy=true")
		//u.Path = "base" // 这边若是赋值了，做转发的时候，会带上path前缀，例： /hello -> /base/hello

		// var query url.Values
		// query.Add("token", "VjouhpQHa6wgWvtkPQeDZbQd")
		// u.RawQuery = query.Encode()

		urls, err := url.Parse(target)
		if err != nil {
			c.String(500, err.Error())
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(urls)
		//proxy := httputil.ReverseProxy{}
		//proxy.Director = func(req *http.Request) {
		//	fmt.Println(req.URL.String())
		//	req.URL.Scheme = "http"
		//	req.URL.Host = "172.16.60.161"
		//	rawQ := req.URL.Query()
		//	rawQ.Add("token", "VjouhpQHa6wgWvtkPQeDZbQd")
		//	req.URL.RawQuery = rawQ.Encode()
		//}

		// proxy.ErrorHandler // 可以添加错误回调
		// proxy.Transport // 若有需要可以自定义 http.Transport
		proxy.ModifyResponse = func(resp *http.Response) error {
			// resp.Header.Del("Server")
			resp.Header.Add("Access-Control-Allow-Origin", "*")
			resp.Header.Add("PROXY", "lflxp-proxy")
			return nil
		}
		proxy.Transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		}
		proxy.Director = func(request *http.Request) {
			targetQuery := urls.RawQuery
			request.URL.Scheme = urls.Scheme
			request.URL.Host = urls.Host
			request.Host = urls.Host // todo 这个是关键
			request.URL.Path = urls.Path
			request.URL.RawPath = urls.RawPath

			if targetQuery == "" || request.URL.RawQuery == "" {
				request.URL.RawQuery = targetQuery + request.URL.RawQuery
			} else {
				request.URL.RawQuery = targetQuery + "&" + request.URL.RawQuery
			}
			if _, ok := request.Header["User-Agent"]; !ok {
				// explicitly disable User-Agent so it's not set to default value
				request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36")
			}
			log.Println("request.URL.Path：", request.URL.Path, "request.URL.RawQuery：", request.URL.RawQuery)
		}
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}

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
	// token := "VjouhpQHa6wgWvtkPQeDZbQd"

	rawUrl.Scheme = proxyUrl.Scheme
	rawUrl.Host = proxyUrl.Host
	// ruq := rawUrl.Query()
	// ruq.Add("token", token)
	// rawUrl.RawQuery = ruq.Encode()
	return nil
}

// 自定义http七层代理服务
func NewHttpProxyByGinCustom(target string, filter map[string]string) func(c *gin.Context) {
	return func(c *gin.Context) {
		if len(filter) > 0 {
			for key, value := range filter {
				if c.GetHeader(key) != value {
					c.String(http.StatusBadRequest, fmt.Sprintf("Header %s:%s not define", key, value))
					return
				}
			}
		}

		err := setTokenToUrl(target, c.Request.URL)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("填写的地址有误: %s", err.Error()))
			c.Abort()
			return
		}

		req, err := http.NewRequestWithContext(c, c.Request.Method, c.Request.URL.String(), c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			c.Abort()
			return
		}
		defer req.Body.Close()
		req.Header = c.Request.Header

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			c.Abort()
			return
		}
		// header 也带过来
		for k := range resp.Header {
			for j := range resp.Header[k] {
				c.Header(k, resp.Header[k][j])
			}
		}
		extraHeaders := make(map[string]string)
		extraHeaders["PROXY"] = "lflxp-proxy"
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, extraHeaders)
		c.Abort()
	}
}

// http 高可用七层代理
func NewHttpProxyHA(target []string, filter map[string]string) func(c *gin.Context) {
	return func(c *gin.Context) {
		if len(filter) > 0 {
			for key, value := range filter {
				if c.GetHeader(key) != value {
					c.String(http.StatusBadRequest, fmt.Sprintf("Header %s:%s not define", key, value))
					return
				}
			}
		}

		if len(target) > 0 {
			urls := []*url.URL{}
			for _, u := range target {
				tmp, err := url.Parse(u)
				if err != nil {
					c.String(http.StatusBadRequest, fmt.Sprintf("target %s is not a valid URL", u))
					return
				}
				urls = append(urls, tmp)
			}

			num := rand.Int() % len(urls)
			director := func(req *http.Request) {
				tg := urls[num]
				targetQuery := tg.RawQuery
				req.URL.Scheme = tg.Scheme
				req.URL.Host = tg.Host
				req.URL.Path = tg.Path
				req.Host = tg.Host // todo 这个是关键
				req.URL.RawPath = tg.RawPath

				if targetQuery == "" || req.URL.RawQuery == "" {
					req.URL.RawQuery = targetQuery + req.URL.RawQuery
				} else {
					req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
				}
				if _, ok := req.Header["User-Agent"]; !ok {
					// explicitly disable User-Agent so it's not set to default value
					req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36")
				}
				log.Println("request.URL.Path：", req.URL.Path, "request.URL.RawQuery：", req.URL.RawQuery)
			}

			proxy := httputil.ReverseProxy{Director: director}
			proxy.ModifyResponse = func(resp *http.Response) error {
				resp.Header.Add("TARGET", urls[num].String())
				return nil
			}
			proxy.Transport = &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				Dial: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 10 * time.Second,
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			}
			proxy.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		} else {
			c.String(404, "target is none")
		}
	}
}
