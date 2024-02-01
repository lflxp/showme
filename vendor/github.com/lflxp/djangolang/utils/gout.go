package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/dataflow"
)

var (
	onceHttpClient                   sync.Once
	httpClient                       *gout.Client
	httpClientWithInsecureSkipVerify *gout.Client
)

// 单例工厂模式
func newGoutClientSingle(skipVerify bool) *gout.Client {
	onceHttpClient.Do(func() {
		httpClient = gout.NewWithOpt()
		httpClientWithInsecureSkipVerify = gout.NewWithOpt(gout.WithInsecureSkipVerify())
	})

	if skipVerify {
		return httpClientWithInsecureSkipVerify
	}
	return httpClient
}

type GoutCli struct {
	client     *gout.Client
	skipVerify bool
}

// 默认开启tls忽略验证
func NewGoutClient() *GoutCli {
	g := &GoutCli{
		skipVerify: true,
	}
	g.client = newGoutClientSingle(true)
	return g
}

// 自定义tls证书校验模式
func (g *GoutCli) SetSkipVerify(skipVerify bool) *GoutCli {
	g.skipVerify = skipVerify
	g.client = newGoutClientSingle(skipVerify)
	return g
}

// 统一处理k8s结构体异常的错误
func (g *GoutCli) GET(url string) *dataflow.DataFlow {
	return g.client.GET(url).ResponseUse(NewCodeGoutResponseUse([]int{http.StatusOK, http.StatusCreated}))
}

func (g *GoutCli) PUT(url string) *dataflow.DataFlow {
	return g.client.PUT(url).ResponseUse(NewCodeGoutResponseUse([]int{http.StatusOK, http.StatusCreated}))
}

func (g *GoutCli) PATCH(url string) *dataflow.DataFlow {
	return g.client.PATCH(url).ResponseUse(NewCodeGoutResponseUse([]int{http.StatusOK, http.StatusCreated}))
}

func (g *GoutCli) POST(url string) *dataflow.DataFlow {
	return g.client.POST(url).ResponseUse(NewCodeGoutResponseUse([]int{http.StatusOK, http.StatusCreated, http.StatusAccepted}))
}

func (g *GoutCli) DELETE(url string) *dataflow.DataFlow {
	return g.client.DELETE(url).ResponseUse(NewCodeGoutResponseUse([]int{http.StatusOK, http.StatusCreated, http.StatusNoContent, http.StatusAccepted}))
}

func (g *GoutCli) OPTIONS(url string) *dataflow.DataFlow {
	return g.client.OPTIONS(url).ResponseUse(NewCodeGoutResponseUse([]int{http.StatusOK, http.StatusCreated}))
}

type GoutError struct {
	Type    string      `json:"type"`
	Status  interface{} `json:"status"`
	Code    interface{} `json:"code"` // code有可能是一个结构体
	Message string      `json:"message"`
}

func (ge GoutError) Error() string {
	codeS := ""
	switch v := ge.Code.(type) {
	case int:
		codeS = strconv.Itoa(v)
	case string:
		codeS = v
	default:
		b, _ := json.Marshal(ge.Code)
		codeS = string(b)
	}

	return fmt.Sprintf("code: %s, message: %s", codeS, ge.Message)
}

type GoutResponseMiddleware struct{}

// ModifyResponse 统一对gout请求的response中的code、date字段处理
func (gout *GoutResponseMiddleware) ModifyResponse(response *http.Response) error {
	use := NewModifyResponse([]int{http.StatusOK, http.StatusCreated})
	return use.ModifyResponse(response)
}

type GoutResponseDelete struct{}

// ModifyResponse 删除时 200、202都为成功
// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#delete-secret-v1-core
func (g *GoutResponseDelete) ModifyResponse(response *http.Response) error {
	use := NewModifyResponse([]int{http.StatusOK, http.StatusAccepted})
	return use.ModifyResponse(response)
}

type ModifyResponseUse struct {
	codes []int
}

func NewModifyResponse(codes []int) *ModifyResponseUse {
	return &ModifyResponseUse{codes: codes}

}
func (c *ModifyResponseUse) ModifyResponse(response *http.Response) error {
	for _, code := range c.codes {
		if response.StatusCode == code {
			return nil
		}
	}
	goutErr := GoutError{}
	bodyStr, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if len(bodyStr) == 0 {
		goutErr.Code = fmt.Sprintf("%d", response.StatusCode)
		goutErr.Message = "no body"
		return goutErr
	}

	err = json.Unmarshal(bodyStr, &goutErr)
	if err != nil {
		return err
	}

	return goutErr
}

// CodeGoutResponseUse 按code处理gout请求的response，方法不太公用，先放在这里
type CodeGoutResponseUse struct {
	codes    []int
	check404 bool // true: statusCode=404时额外判断是GET请求不报错
}

func NewCodeGoutResponseUse(codes []int, check ...bool) *CodeGoutResponseUse {
	r := &CodeGoutResponseUse{
		codes: codes,
	}
	if len(check) != 0 {
		r.check404 = check[0]
	} else {
		r.check404 = true
	}
	return r
}
func (c *CodeGoutResponseUse) ModifyResponse(response *http.Response) error {
	goutErr := GoutError{}

	for _, code := range c.codes {
		if response.StatusCode == code {
			return nil
		}
	}
	if c.check404 && response.StatusCode == 404 {
		if strings.ToLower(response.Request.Method) == "get" {
			goutErr.Code = fmt.Sprintf("%d", response.StatusCode)
			goutErr.Message = "资源不存在"
			return goutErr
		}

		bodyStr, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		if len(bodyStr) == 0 {
			return fmt.Errorf("code(%d) no body", response.StatusCode)
		}

		goutErr.Code = fmt.Sprintf("%d", response.StatusCode)
		goutErr.Message = string(bodyStr)
		return goutErr
	}

	bodyStr, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if len(bodyStr) == 0 {
		return fmt.Errorf("code(%d) no body", response.StatusCode)
	}

	err = json.Unmarshal(bodyStr, &goutErr)
	if err != nil {
		goutErr.Code = fmt.Sprintf("%d", response.StatusCode)
		goutErr.Message = string(bodyStr)
		return goutErr
	}

	return goutErr
}
