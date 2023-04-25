package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	Message   chan string
}

// ListenMessage 广播队列，发送函数
func (this *Server) ListenMessage() {
	for {
		fmt.Println("BroadCast channel waiting working...")
		msg := <-this.Message //server channel中取出message,没数据就阻塞
		fmt.Println("BroadCast channel get msg...")
		this.mapLock.Lock()
		//需要广播的数据放入每一个用户的chan
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// BroadCast 加入广播队列
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
	log.Printf(msg)

}

func (this *Server) Handler(conn net.Conn) {
	fmt.Println("connect ok ,handel data:")
	//user := NewUser(conn)
	user := NewUser(conn, this)
	user.Online()
	isLive := make(chan bool)

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if err != nil && err != io.EOF {
				fmt.Println("socket read err", err)
				return
			}
			//msg 处理消息
			user.DoMessage(string(buf[:n-1]))
			isLive <- true

		}
	}()

	for {
		select {
		case <-isLive:
		case <-time.After(time.Second * 600):
			log.Printf(user.Name + " timeout, force offline")
			offlineMsg := "[im_server] your are leave room too long.. offline.\n"
			user.SendMsg(offlineMsg, user)
			user.Offline()
			return
		}
	}
}

// Start this 是server的抽象
func (this *Server) Start() {
	pid := os.Getpid()
	log.Printf("-----imserver run in pid:%d ok-----", pid)
	listener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("net.Listen err", err)
		return
	}
	defer listener.Close() //return 之前执行
	fmt.Println("net.Listen ok waiting for connect..")
	//广播队列中 协程启动
	go this.ListenMessage()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("net.Accept err", err)
			return
		}
		go this.Handler(conn)
	}
}
