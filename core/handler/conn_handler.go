package handler

import (
	"baozun.com/security-proxy/common/log"
	"baozun.com/security-proxy/common/utils"
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

const (
	DefaultDialTimeout = 3 * time.Second
)

type ConnHandler struct {
	Conn net.Conn
}

func (c *ConnHandler) Do() error {
	cReader := NewCacheReader(c.Conn)
	bReader := bufio.NewReader(cReader)
	req, err := http.ReadRequest(bReader)
	if nil != err {
		return fmt.Errorf("http.ReadRequest err:%w", err)
	}
	serverAddr := utils.GetServerAddress(req)

	// 代理转发
	if req.Method == "CONNECT" {
		// 客户端说要建立连接，代理服务器要回应建立好了，然后才可以像HTTP一样请求访问
		log.Info(fmt.Sprintf("Accepting CONNECT to %s", serverAddr))
		return c.transfer(serverAddr, nil, []byte("HTTP/1.0 200 OK\r\n\r\n"))
	}
	if !req.URL.IsAbs() { // 拦截非代理请求
		log.Error(fmt.Sprintf("invalid request, url not abs:%s", req.URL.String()))
		return c.replyError("HTTP/1.1 500 This is a proxy server. Does not respond to non-proxy requests.\r\n\r\n")
	}
	log.Debug(fmt.Sprintf("recv request %v %v %v %v", req.URL.Path, req.Host, req.Method, req.URL.String()))
	return c.transfer(serverAddr, cReader.Cache(), nil)
}

// 请求转发
func (c *ConnHandler) transfer(addr string, tranData, replyData []byte) error {
	targetConn, err := net.DialTimeout("tcp", addr, DefaultDialTimeout) // 建立tcp 连接拨号
	if err != nil {
		log.Error(fmt.Sprintf("at transfer dial %s fail,error: %s", addr, err.Error()))
		return c.replyError("HTTP/1.1 502 Bad Gateway\r\n\r\n")
	}
	if nil != tranData { // 转发
		_, _ = targetConn.Write(tranData)
	}
	if nil != replyData { // 应答
		_, _ = c.Conn.Write(replyData)
	}
	c.copy(targetConn)
	return nil
}

// 应答客户端请求错误信息, 并关闭客户端连接
func (c *ConnHandler) replyError(response string) error {
	if _, err := io.WriteString(c.Conn, response); err != nil {
		return fmt.Errorf("at respError error responding to client: %w", err)
	}
	if err := c.Conn.Close(); err != nil {
		return fmt.Errorf("at respError error closing client connection: %w", err)
	}
	return nil
}

// 请求转发
func (c *ConnHandler) copy(targetConn net.Conn) {
	defer c.Conn.Close()
	go io.Copy(targetConn, c.Conn)
	io.Copy(c.Conn, targetConn)
}
