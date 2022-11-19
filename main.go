package main

import (
	"fmt"
	"github.com/go-mogu/hz-framework/bootstrap"
	"github.com/go-mogu/hz-framework/config"
	"github.com/go-mogu/hz-framework/router"
	"github.com/urfave/cli/v2"
	"os"
	"runtime"
	"time"
)

var (
	// AppName 当前应用名称
	AppName  = "mogu-sms"
	AppUsage = "使用hertz框架作为基础开发库，封装一套适用于面向api编程的快速开发框架"
	// AuthorName 作者
	AuthorName  = "DingDing"
	AuthorEmail = "15077731547@163.com"
	//	AppPort 程序启动端口
	AppPort string
	//	BuildVersion 编译的app版本
	BuildVersion string
	//	BuildAt 编译时间
	BuildAt string
	_UI     = `
 ████     ████   ███████     ████████  ██     ██        ████████ ████     ████  ████████
░██░██   ██░██  ██░░░░░██   ██░░░░░░██░██    ░██       ██░░░░░░ ░██░██   ██░██ ██░░░░░░ 
░██░░██ ██ ░██ ██     ░░██ ██      ░░ ░██    ░██      ░██       ░██░░██ ██ ░██░██       
░██ ░░███  ░██░██      ░██░██         ░██    ░██ █████░█████████░██ ░░███  ░██░█████████
░██  ░░█   ░██░██      ░██░██    █████░██    ░██░░░░░ ░░░░░░░░██░██  ░░█   ░██░░░░░░░░██
░██   ░    ░██░░██     ██ ░░██  ░░░░██░██    ░██             ░██░██   ░    ░██       ░██
░██        ░██ ░░███████   ░░████████ ░░███████        ████████ ░██        ░██ ████████ 
░░         ░░   ░░░░░░░     ░░░░░░░░   ░░░░░░░        ░░░░░░░░  ░░         ░░ ░░░░░░░░  

`
)

// Stack 程序运行前的处理
func Stack() *cli.App {
	buildInfo := fmt.Sprintf("%s-%s-%s-%s-%s", runtime.GOOS, runtime.GOARCH, BuildVersion, BuildAt, time.Now().Format("2006-01-02 15:04:05"))

	return &cli.App{
		Name:    AppName,
		Version: buildInfo,
		Usage:   AppUsage,
		Authors: []*cli.Author{
			{
				Name:  AuthorName,
				Email: AuthorEmail,
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "env",
				Aliases:     []string{"e"},
				Value:       "dev",
				Usage:       "请选择配置文件 [dev | test | prod]",
				Destination: &config.ConfEnv,
			},
			&cli.StringFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Value:       "9604",
				Usage:       "请选择启动端口",
				Destination: &AppPort,
			},
		},
		Action: func(context *cli.Context) error {
			fmt.Println(fmt.Sprintf("\u001B[34m%s\u001B[0m", _UI))

			//	程序启动时需要加载的服务
			bootstrap.BootService()
			//	注册路由 启动程序
			router.Register(AppPort).Spin()
			return nil
		},
	}
}

func main() {
	if err := Stack().Run(os.Args); err != nil {
		panic(err)
	}
}
