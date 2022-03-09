package proxy

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var target string

func customProxy(c *gin.Context) {
	// if c.GetHeader("direct") != "lab" {
	// 	return
	// }

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
	extraHeaders["direct"] = "lab"
	c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, extraHeaders)
	c.Abort()
}

func setTokenToUrl(target string, rawUrl *url.URL) error {
	// 这边是从设置里拿代理值
	//equipment, err := proxy.GetEquipment()
	//if err != nil {
	//	return err
	//}

	if !strings.Contains(target, "http") {
		target = fmt.Sprintf("http://%s", target)
	}
	proxyUrl := target
	// token := "VjouhpQHa6wgWvtkPQeDZbQd"
	u, err := url.Parse(proxyUrl)
	if err != nil {
		return err
	}

	rawUrl.Scheme = u.Scheme
	rawUrl.Host = u.Host
	// ruq := rawUrl.Query()
	// ruq.Add("token", token)
	// rawUrl.RawQuery = ruq.Encode()
	return nil
}

func mutilHAProxy(c *gin.Context) {
	one, err := url.Parse("https://localhost")
	if err != nil {
		c.String(500, err.Error())
		return
	}
	urls := []*url.URL{
		one,
		{
			Scheme: "https",
			Host:   "localhost",
		},
	}

	num := rand.Int() % len(urls)
	director := func(req *http.Request) {
		target := urls[num]
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
	}

	proxy := httputil.ReverseProxy{Director: director}
	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Del("Server")
		resp.Header.Add("urlsId", fmt.Sprintf("%d", num))
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
}

func Run(sourcePort, targetUrl string) {
	target = targetUrl
	r := gin.Default()

	r.Any("/*action", customProxy)
	r.Run(fmt.Sprintf(":%s", sourcePort))
}
