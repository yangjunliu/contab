package main

import (
	"flag"
	"fmt"
	"github.com/yangjunliu/crontab/master"
	"runtime"
	"time"
)

var (
	confFile string
)

func initArgs() {
	flag.StringVar(&confFile, "config", "./master.json", "传入master配置")
	flag.Parse()
}

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	// 初始化命令行参数
	initArgs()

	// 初始化线程数
	initEnv()

	// 加载配置
	if err := master.InitConfig(confFile); err != nil {
		fmt.Println(err)
	}

	// 启动服务
	if err := master.InitApiServer(); err != nil {
		fmt.Println(err)
	}

	// 任务管理器
	if err := master.InitJobMgr(); err != nil {
		fmt.Println(err)
	}

	for {
		time.Sleep(5 * time.Second)
	}
}
