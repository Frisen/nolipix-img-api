package response

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code" example:"0"`      // 状态码 0:成功 1:异常
	Data interface{} `json:"data"`                  // 数据
	Msg  string      `json:"msg" example:"success"` // 消息
}

type Page struct {
	List      interface{} `json:"list"`                   // 数据列表
	Count     int         `json:"count" example:"50"`     // 数据总数
	PageIndex int         `json:"page_index" example:"1"` // 当前页码
	PageSize  int         `json:"page_size" example:"20"` // 当前页数据总数
	Extend    interface{} `json:"extend"`                 //扩展字段
}

type PageAll struct {
	List interface{} `json:"list"` // 数据列表
}

type ResponseOk struct {
	Code int         `json:"code"` // 状态码 0:成功 1:异常
	Data interface{} `json:"data"` // 数据
	Msg  string      `json:"msg"`  // 消息
}

func OK(c *gin.Context, data interface{}, msg string) {
	response := ResponseOk{
		Code: 0,
		Data: data,
		Msg:  msg,
	}
	write(c, response)
}

func SimpleOK(c *gin.Context) {
	OK(c, nil, "")
}

func Error(c *gin.Context, data interface{}, msg string) {
	response := Response{
		Data: data,
		Msg:  msg,
		Code: 1,
	}
	c.JSON(http.StatusOK, response)
}

func PageAllOK(c *gin.Context, result interface{}, msg string) {

	response := Response{
		Data: PageAll{
			List: result,
		},
		Msg:  msg,
		Code: 0,
	}

	c.JSON(http.StatusOK, response)
}

func write(c *gin.Context, resp interface{}) {
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Header().Set("content-type", "application/json;charset=utf-8")

	//防止转义html
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(resp)

	_, _ = c.Writer.Write(buffer.Bytes())
}
