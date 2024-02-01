package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

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
	ErrorMessage string      `json:"message"`
	Host         string      `json:"host"`
	TraceId      string      `json:"traceid"`
	ShowType     string      `json:"showtype"`
	Code         int         `json:"code"`
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
		Success:      true,
		Data:         data,
		Host:         c.Request.URL.Path,
		Code:         20000,
		ErrorMessage: "success",
	}

	c.JSONP(code, info)
}

func SendErrorMessageInterface(c *gin.Context, code int, codeStr string, data interface{}) {
	info := &Result{
		Success:      true,
		Data:         data,
		Host:         c.Request.URL.Path,
		Code:         0,
		ErrorMessage: codeStr,
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
		Code:         0,
		ErrorMessage: errorMsg,
		Host:         c.Request.URL.Path,
	}

	c.Writer.Header().Del("Content-Encoding")
	c.JSONP(code, info)
}

// SendError...
func SendError(c *gin.Context, code int, err error) {
	errorCode := FailedUnknown
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
		switch err.(type) {
		case GoutError:
			if err.(GoutError).Code == fmt.Sprintf("%d", http.StatusNotFound) {
				errorCode = ResourceNotFound
				errorMsg = err.(GoutError).Message
			}
		}
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

// PkgResp 包装返回类
func PkgResp(c *gin.Context, code int, result Result) {
	result.Host = c.Request.URL.Path
	c.JSONP(code, result)
}

func CliErr(err error) Result {
	slog.Error("get rancher v3-default client error", "Error", err.Error())
	return Result{
		Data:         err.Error(),
		Success:      false,
		ErrorCode:    FailedRemoteServiceError,
		ErrorMessage: fmt.Sprintf("client加载失败 %s", err.Error()),
	}
}

func PkgSuccess(data interface{}) Result {
	return Result{
		Data:    data,
		Success: true,
	}
}

func PkgErrorMsg(errMsg string) Result {
	return Result{
		Success:      false,
		ErrorCode:    strconv.Itoa(http.StatusInternalServerError),
		ErrorMessage: fmt.Sprintf(errMsg),
	}
}

// BadDecode 解码失败
func BadDecode(ctx *gin.Context, err error) {
	ctx.JSONP(http.StatusBadRequest, Result{
		Data:         err.Error(),
		Success:      false,
		ErrorCode:    FailedDecodeError,
		ErrorMessage: fmt.Sprintf("Rsa解码失败 %s", err.Error()),
		Host:         ctx.Request.URL.Path,
	})
}

// BadParams 参数绑定失败
func BadParams(ctx *gin.Context, err error) {
	ctx.JSONP(http.StatusBadRequest, Result{
		Data:         err.Error(),
		Success:      false,
		ErrorCode:    FailedParamsError,
		ErrorMessage: fmt.Sprintf("参数绑定失败 %s", err.Error()),
		Host:         ctx.Request.URL.Path,
	})
}

type E422 struct {
	Type    string      `json:"type"`
	Links   interface{} `json:"links"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Status  int         `json:"status"`
}

func (e *E422) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func ParseE422(data []byte) (*E422, error) {
	var result *E422
	err := json.Unmarshal(data, result)
	return result, err
}

// -----------------v3用的错误返回信息------------------
var rancherErrorMap = map[int]string{
	401: "账号权限不足",
	409: "提交内容中可能包含两个冲突版本，请刷新后重试",
	404: "请求资源不存在",
}

// 硬编码翻译部分错误
func translateErrorMsg(errMsg string) string {
	if strings.Contains(errMsg, "failed to decode key: not valid pem format") {
		return "证书或私钥格式错误，请检查"
	}
	if strings.Contains(errMsg, "mountPath: Invalid value") {
		if strings.Contains(errMsg, "must be unique") {
			// TODO 定位比较丑，需优化
			startIndex := strings.Index(errMsg, "mountPath: Invalid value")
			endIndex := strings.Index(errMsg, "must be unique")
			return fmt.Sprintf("挂载路径 %s 重复", errMsg[startIndex+24:endIndex])
		}
	}
	return errMsg
}

// -----------------错误返回解析END------------------
// -----------------v3用的错误返回信息END------------------
