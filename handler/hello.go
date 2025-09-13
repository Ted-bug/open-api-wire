package handler

import (
	"api-gin/repo"
	"github.com/gin-gonic/gin"
)

type HelloHandler struct {
	helloRepo *repo.UserRepo
}

func NewHelloHandler(
	helloRepo *repo.UserRepo,
) *HelloHandler {
	return &HelloHandler{
		helloRepo: helloRepo,
	}
}

func (h *HelloHandler) Hello(ctx *gin.Context, req *HelloReq) (resp *HelloResp, err error) {
	return &HelloResp{
		Msg: h.helloRepo.Hello(ctx) + req.Name,
	}, nil
}
