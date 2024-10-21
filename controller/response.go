package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
	{
		"code":
		"msg":
		"data":
	}
*/
type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy
	CodeNeedLogin
	CodeInvalidToken
)

var getMsg = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户已存在",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",
	CodeNeedLogin:       "需要登录",
	CodeInvalidToken:    "无效的token",
}

type Response struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"` // 为空忽略
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	rd := &Response{
		Code: CodeSuccess,
		Msg:  getMsg[CodeSuccess],
		Data: data,
	}
	c.JSON(http.StatusOK, rd)
}

func ResponseError(c *gin.Context, code ResCode) {
	rd := &Response{
		Code: code,
		Msg:  getMsg[code],
		Data: nil,
	}
	c.JSON(http.StatusOK, rd)
}

func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	rd := &Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
	c.JSON(http.StatusOK, rd)
}
