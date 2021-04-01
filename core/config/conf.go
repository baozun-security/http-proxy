package config

import (
	"baozun.com/security-proxy/common/log"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
)

var (
	CoreConf *Settings
)

func Init(conf string) {
	_, err := toml.DecodeFile(conf, &CoreConf)
	if err != nil {
		fmt.Printf("Err %v", err)
		os.Exit(1)
	}
}

type Settings struct {
	Log    log.Config
	Server Server
}

type Server struct {
	Port        int
	MaxWorkers  int // 协程工作者数
	MaxJobQueue int // 工作队列大小
}
