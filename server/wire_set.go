package server

import (
	"api-gin/config"
	"api-gin/controller"
	"api-gin/handler"
	"api-gin/infra/log"
	"api-gin/repo"
	"github.com/google/wire"
)

var (
	baseSet = wire.NewSet(
		config.NewConfig,
		log.NewLogger,
	)
	repoSet = wire.NewSet(
		repo.NewDB,
		repo.NewUserRepo,
	)
	// serviceSet = wire.NewSet()
	handlerSet = wire.NewSet(
		handler.NewHelloHandler,
	)
	controllerSet = wire.NewSet(
		controller.NewHelloController,
	)
)

type Controllers struct {
	// 加入控制器层
	HelloController *controller.HelloController
}
