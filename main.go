package main

import (
	"baozun.com/security-proxy/common/log"
	"baozun.com/security-proxy/core/config"
	"baozun.com/security-proxy/core/server"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"os"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

/*正向代理服务
1: 支持http、https流量代理，以及分析统计
2: 转发流量到后端xray(安全扫描工具)
*/
func main() {
	var (
		cfg string
	)
	flag.StringVar(&cfg, "conf", "", "server config [toml]")
	flag.Parse()
	if len(cfg) == 0 {
		fmt.Println("config is empty")
		os.Exit(0)
	}
	config.Init(cfg) // 配置文件加载
	conf := config.CoreConf
	log.Init(&conf.Log)                                                        // 初始化日志
	server.Init(conf.Server.MaxWorkers, conf.Server.MaxJobQueue, log.Logger()) // 初始化服务
	if err := server.Run(); nil != err {
		log.Fatal("server run error", zap.Error(err))
	}
}
