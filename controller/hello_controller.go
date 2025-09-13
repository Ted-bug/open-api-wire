package controller

import (
	"api-gin/handler"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HelloController struct {
	h *handler.HelloHandler
}

func NewHelloController(h *handler.HelloHandler) *HelloController {
	return &HelloController{
		h: h,
	}
}

func (c *HelloController) Hello(ctx *gin.Context) {
	var req handler.HelloReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	resp, err := c.h.Hello(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "内部错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}
