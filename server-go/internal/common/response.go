package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	SuccessCode = iota
	ErrorCode

	AuthFailedErrorCode
	DuplicateFieldErrorCode

	RedisOperateErrorCode
	ParamInvalidErrorCode

	TemporarilyUnavailable
	ErrorIDLen
	ErrorID
	UnknownErrorCode
	TimeOutErrorCode
	RemoteAllFailedErrorCode
	ProjectIDRespect
	SomeFieldIsNull
	ExceedLimitErrorCode
	SSO_ERROR
	OTPErrorCode
	PasswordNeedChanged
	NeedCaptchaCheck
	LoginIpNotInWhiteList
	UserLocked

	//业务错误码
	DirISNotExists
	FileIsNotExists
	ReadDirFailed
)

var ErrorDescriptions = map[int]string{
	SuccessCode:              "success",
	AuthFailedErrorCode:      "auth failed",
	DuplicateFieldErrorCode:  "duplicate field",
	RedisOperateErrorCode:    "redis operate error",
	ParamInvalidErrorCode:    "param invalid",
	TemporarilyUnavailable:   "resource temporarily unavailable",
	ErrorIDLen:               "ID MAX LEN IS 1-15",
	ErrorID:                  "ID ONLY SYUUPRT 'A-Z/a-z/0-9/-/_'",
	UnknownErrorCode:         "unknown error",
	ProjectIDRespect:         "PROJECT ID REPECT",
	SomeFieldIsNull:          "SOME FIELD IS NUL",
	TimeOutErrorCode:         "get result timeout",
	RemoteAllFailedErrorCode: "all remote instance failed",
	SSO_ERROR:                "sso error",
	OTPErrorCode:             "otp required",
	PasswordNeedChanged:      "password has not been updated for a long time",
	NeedCaptchaCheck:         "need captcha check",
	LoginIpNotInWhiteList:    "login ip not in whitelist",
	UserLocked:               "user locked",

	//业务错误码
	DirISNotExists:  "did is not exists",
	FileIsNotExists: "file is not exists",
	ReadDirFailed:   "read dir failed",
}

type Response struct {
	ErrorCode int         `json:"errorCode"`
	Message   string      `json:"msg"`
	Data      interface{} `json:"data"`
}

func (response *Response) SetError(code int) {
	response.ErrorCode = code

	if msg, ok := ErrorDescriptions[code]; ok {
		response.Message = msg
	}
}

func CreateResponse(c *gin.Context, code int, data interface{}) {
	var response Response

	response.SetError(code)
	response.Data = data
	c.JSON(
		http.StatusOK,
		response,
	)
}

func BadRequest(c *gin.Context, msg string) {
	var response Response

	response.SetError(ParamInvalidErrorCode)
	response.Message = msg
	c.JSON(
		http.StatusBadRequest,
		response,
	)
}

func Unauthorized(c *gin.Context, msg string) {
	var response Response

	response.SetError(AuthFailedErrorCode)
	response.Message = msg
	c.JSON(
		http.StatusUnauthorized,
		response,
	)
}

// 创建符合 Ant Design 前端要求的响应。这些函数将数据转换为 AntResponseData 结构，并将其作为 JSON 响应发送给前端
type AntResponseData struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}

func CreateAntResponse(c *gin.Context, code int, data interface{}) {
	var response AntResponseData
	response.Data = data
	response.Success = (code == SuccessCode)

	c.JSON(http.StatusOK, response)
}

func CreateAntErrorResponse(c *gin.Context, code int, msg string) {
	var response AntResponseData
	response.Data = nil
	response.Success = false

	c.JSON(http.StatusOK, response)
}
