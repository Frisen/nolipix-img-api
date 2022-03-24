package api

import (
	"nolipix-img-api/api/dto"
	"nolipix-img-api/api/response"
	"nolipix-img-api/service"

	"github.com/gin-gonic/gin"
)

func Compress(c *gin.Context) {
	req := new(dto.CompressReq)
	err := c.ShouldBindJSON(req)
	if err != nil {
		response.Error(c, req, err.Error())
		return
	}
	thumbUrl, err := service.Compress(req.Url, req.Rows, req.Cols)
	if err != nil {
		response.Error(c, thumbUrl, err.Error())
		return
	}
	response.OK(c, thumbUrl, "")
	return
}
