package cmd

import (
	"api-gin/server"
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start server",
	Long:  `start server`,
	Run:   startCmdExculpate,
}

func startCmdExculpate(cmd *cobra.Command, args []string) {
	app, err := server.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	addr := fmt.Sprintf("%s:%d", app.Host, app.Port)
	fmt.Printf("点击访问: https://%s\n", addr)
	srv := &http.Server{
		Addr:    addr,
		Handler: app.Engine,
	}
	// 4.1 开启一个goroutine处理请求；否则会一直循环中，无法执行往下的关闭代码
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen error: %s\n", err)
		}
	}()
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	pid := os.Getpid()
	_ = os.WriteFile(wd+"/app.pid", []byte(fmt.Sprintf("%d", pid)), 0644)

	// 4.2 创建一个通道监听中断信号
	// kill（syscall.SIGTERM）、kill -2(syscall.SIGINT)监听得到、kill -9监听不到
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit // 无信号会阻塞
	fmt.Println("shutdown server ...")
	// 4.3 接收到结束信号，创建5秒超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown: ", err)
	}
	fmt.Println("server exiting")
}
