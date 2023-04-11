package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}
func (this *Server) Handler(conn net.Conn) {
	fmt.Println("connect ok ,handel data..")
}

//this 是server的抽象
func (this *Server) Start() {
	listener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("net.Listen err", err)
		return
	}
	defer listener.Close()
	fmt.Println("net.Listen ok waiting for connect..")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("net.Accept err", err)
			return
		}
		go this.Handler(conn)
	}

}
