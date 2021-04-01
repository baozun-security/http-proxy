package server

import (
	"baozun.com/security-proxy/common/async"
	"baozun.com/security-proxy/common/log"
	"baozun.com/security-proxy/core/config"
	"baozun.com/security-proxy/core/handler"
	"fmt"
	"go.uber.org/zap"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	quit chan struct{}
	pool *async.WorkerPool
)

// Init proxy server
func Init(maxWorkers, maxQueue int, log *zap.Logger) {
	pool = async.NewWorkerPool(maxWorkers, maxQueue, log).Run() // init worker pool
}

// Run start tcp listen server
func Run() error {
	addr := fmt.Sprintf(":%d", config.CoreConf.Server.Port)
	lis, err := net.Listen("tcp", addr)
	if nil != err {
		return fmt.Errorf("failed get listen addr, err:%s", err.Error())
	}
	go tryDisConn()       // 优雅退出
	return listening(lis) // 持续监听TCP连接请求
}

// 启动监听
func listening(lis net.Listener) error {
	for {
		select {
		case <-quit:
			log.Info("tcp listening quited")
			return nil
		default: // 阻塞监听连接请求
			conn, err := lis.Accept()
			if nil != err {
				return fmt.Errorf("listen accept err:%s", err.Error())
			}
			pool.Add(&handler.ConnHandler{Conn: conn})
		}
	}
}

func tryDisConn() {
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGKILL,
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP,
		syscall.SIGABRT,
	)

	select {
	case sig := <-signals:
		go func() { // 10s之后强制退出系统
			select {
			case <-time.After(time.Second * 10):
				log.Warn("Shutdown gracefully timeout, application will shutdown immediately.")
				os.Exit(0)
			}
		}()
		log.Info(fmt.Sprintf("get signal %s, application will shutdown.", sig))

		log.Debug("Start Stop ProxyServer")
		quit <- struct{}{}

		os.Exit(0)
	}
}
