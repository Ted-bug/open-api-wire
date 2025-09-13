package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var stoptCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop server",
	Long:  `stop server`,
	Run:   stopCmdExculpate,
}

func stopCmdExculpate(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("工作路径错误：%v", err)
	}
	pid, err := os.ReadFile(wd + "/app.pid")
	if err != nil {
		log.Fatalf("读取进程号错误：%v", err)
	}
	pidStr := strings.Trim(string(pid), "\n")
	if pidStr == "" {
		log.Fatalf("进程号错误")
	}
	pidInt, err := strconv.Atoi(pidStr)
	if err != nil {
		log.Fatalf("进程号错误：%v", err)
	}
	fmt.Println(pidInt)
	process, err := os.FindProcess(pidInt)
	if err != nil {
		log.Fatalf("查找进程错误：%v", err)
	}
	defer func() {
		_ = process.Release()
	}()
	doneError := make(chan struct{})
	waitError := make(chan struct{})
	defer func() {
		close(doneError)
		close(waitError)
	}()
	go func() {
		err = process.Signal(os.Interrupt)
		if err != nil {
			waitError <- struct{}{}
		}
	}()
	go func() {
		for {
			err := process.Signal(syscall.Signal(0))
			if err != nil {
				doneError <- struct{}{}
			}
			time.Sleep(20 * time.Millisecond)
		}
	}()
	select {
	case <-doneError:
		log.Println("服务已关闭")
	case <-waitError:
		log.Fatalf("进程号错误：%v", err)
	case <-time.After(time.Second * 5):
		log.Fatalf("关闭等待超时")
	}
}
