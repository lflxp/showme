package proxy

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	utils "github.com/lflxp/tools/httpclient"

	"github.com/gin-gonic/gin"
)

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

	result := []string{"", tpaths[len(tpaths)-1]}
	result = append(result, rawpaths[2:]...)

	rawUrl.Path = strings.Join(result, "/")
	rawUrl.Scheme = proxyUrl.Scheme
	rawUrl.Host = proxyUrl.Host
	// ruq := rawUrl.Query()
	// ruq.Add("token", token)
	// rawUrl.RawQuery = ruq.Encode()
	return nil
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

		// check the proxy request whether it is websocket
		// if c.Request.Header.Get("Connection") == "Upgrade" && c.Request.Header.Get("Upgrade") == "websocket" {
		// 	WebSocketProxy(target, c)
		// 	return
		// }

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
		if c.Request.Header.Get("token") != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Request.Header.Get("token")))
		}

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
		extraHeaders["PROXY"] = "proxy"

		// header 也带过来
		for k := range resp.Header {
			for j := range resp.Header[k] {
				c.Header(k, resp.Header[k][j])
			}
		}

		for _, cookie := range resp.Cookies() {
			c.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, c.Request.Host, cookie.Secure, cookie.HttpOnly)
		}

		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, extraHeaders)
		// utils.SendSuccessMessage(c, resp.StatusCode, msg)
		c.Abort()

	}
}
