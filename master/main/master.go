package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/yangjunliu/crontab/master"
)

var (
	confFile string
	pwd      string
	err      error
)

func initArgs() {
	pwd, err = os.Getwd()
	if err != nil {
		pwd = "./master.json"
	} else {
		pwd += "\\master\\main\\master.json"
	}
	log.Println(pwd)

	flag.StringVar(&confFile, "config", pwd, "传入master配置")
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

	// 任务管理器
	if err := master.InitJobMgr(); err != nil {
		fmt.Println(err)
	}

	// 启动服务
	if err := master.InitApiServer(); err != nil {
		fmt.Println(err)
	}
}
