package server

import (
	"api-gin/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type App struct {
	Host        string
	Port        int
	Engine      *gin.Engine  // 引擎
	Controllers *Controllers // router配置
}

func NewApp(config *config.Config, controllers *Controllers) (*App, error) {
	if config == nil {
		return nil, fmt.Errorf("[App] 配置不能为空")
	}
	gin.SetMode(config.Mode)

	g := gin.New()

	loc, _ := time.LoadLocation("Asia/Shanghai")
	time.Local = loc

	// 引入一些中间件
	// g.Use(middleware.RecoveryMiddlerware())
	// g.Use(middleware.LoggerMiddlerware())

	app := &App{
		Host:        config.Host,
		Port:        config.Port,
		Engine:      g,
		Controllers: controllers,
	}
	app.initRouter()

	return app, nil
}

func (a *App) initRouter() {
	rGroup := a.Engine.RouterGroup
	api := rGroup.Group("/v1/api/hello")
	{
		api.GET("/:name", a.Controllers.HelloController.Hello)
	}
}
