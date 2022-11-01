package httpclient

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

// request fail
const (
	Failed                   = "4000"
	FailedParamsError        = "4002"
	FailedRemoteServiceError = "4001"
	FailedUnknown            = "4003"
	FailedDecodeError        = "4004"
	AthorizationError        = "4005"
	JsonError                = "4006"
	ResourceNotFound         = "4007"
	LicenseError             = "4008"
	BiddenError              = "4009"
	RoleChangeNeedReLogin    = "4010"
)

// system error
const (
	SystemError = "9000"
)

var responseMap = map[string]string{
	Failed:                   "操作失败!",
	FailedParamsError:        "参数错误!",
	FailedUnknown:            "操作失败，未知错误!",
	FailedRemoteServiceError: "调用第三方接口异常!",
	SystemError:              "系统异常!",
	FailedDecodeError:        "解码失败",
	AthorizationError:        "Header Not Contains Authorization",
	JsonError:                "json序列化错误",
}

type Result struct {
	Success      bool        `json:"success"`
	Data         interface{} `json:"data"`
	ErrorCode    string      `json:"errorCode"`
	ErrorMessage string      `json:"errorMessage"`
	Host         string      `json:"host"`
	TraceId      string      `json:"traceid"`
	ShowType     string      `json:"showtype"`
}

func (r *Result) String() string {
	str, err := json.Marshal(r)
	if err != nil {
		return err.Error()
	}
	return string(str)
}

func SendMessage(c *gin.Context, code int, success bool, data interface{}, errcode, errmsg, host, traceid, showtype string) {
	info := &Result{
		Success:      success,
		Data:         data,
		ErrorCode:    errcode,
		ErrorMessage: errmsg,
		Host:         host,
		TraceId:      traceid,
		ShowType:     showtype,
	}
	c.JSONP(code, info)
}

// SendSuccessMessage select success and have data
func SendSuccessMessage(c *gin.Context, code int, data interface{}) {
	info := &Result{
		Success: true,
		Data:    data,
		Host:    c.Request.URL.Path,
	}

	c.JSONP(code, info)
}

// SendErrorMessage select fail
func SendErrorMessage(c *gin.Context, code int, errorCode string, errorMsg string) {
	if errorMsg == "" {
		errorMsg = responseMap[errorCode]
	}

	info := &Result{
		Success:      false,
		Data:         nil,
		ErrorCode:    errorCode,
		ErrorMessage: errorMsg,
		Host:         c.Request.URL.Path,
	}

	c.Writer.Header().Del("Content-Encoding")
	c.JSONP(code, info)
}
