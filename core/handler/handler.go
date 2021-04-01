package handler

import "net"

type ConnHandler struct {
	Conn net.Conn
}

func (h *ConnHandler) Do() error {

}