package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server} //传入servers实例

	go user.ListenMessage()

	return user

}

func (this *User) Online() {
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	//
	this.server.BroadCast(this, "online now\n")
	helpList := "=============Welcome to chat room============\n" +
		"off: offline your client\n" +
		"who: list all client\n" +
		"all: send a broadcast message. exp all:hello ==> send 'hello'\n" +
		"rename: change your name.exp rename:xxx ===> your name is xxx\n" +
		"sendto| send message to someone. exp sendto|username|msg\n" +
		"===========================================\n"
	this.SendMsg(helpList, this)

}
func (this *User) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
	//
	this.server.BroadCast(this, "Offline now\n")
	err := this.conn.Close()
	if err != nil {
		defer this.conn.Close()
		return
	}

}

func (this *User) DoMessage(msg string) {

	chatRecord := "chatRecord: " + "[" + this.Addr + "]" + this.Name + ": " + msg + "\n"
	log.Printf(chatRecord)
	if msg == "off" {
		this.Offline()

	} else if msg == "who" {
		// 在线列表
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "online...\n"
			this.SendMsg(onlineMsg, this)
		}
		this.server.mapLock.Unlock()
	} else if msg == "help" {
		// 呼出帮助面板
		helpList := "off: offline your client\n" +
			"who: list all client\n" +
			"all: send a broadcast message. exp all:hello ==> send 'hello'\n" +
			"rename: change your name.exp rename:xxx ===> your name is xxx\n" +
			"sendto| send message to someone. exp sendto|username|msg\n"
		this.SendMsg(helpList, this)
	} else if strings.HasPrefix(msg, "all:") {
		// 广播消息
		this.server.BroadCast(this, msg[4:]+"\n")
	} else if strings.HasPrefix(msg, "rename:") {
		// 修改用户名
		newName := msg[7:]
		log.Printf("rename from " + this.Name + " to " + newName)
		_, used := this.server.OnlineMap[newName]
		if used {
			this.SendMsg("[im_server] name been used,pleas change\n", this)
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()
			this.Name = newName
			this.SendMsg("[im_server] "+"name changed,your name: "+newName+"\n", this)

		}
		//sendto:username:msg
	} else if strings.HasPrefix(msg, "sendto|") {
		remoteName := strings.Split(msg, "|")[1]
		remoteMsg := strings.Split(msg, "|")[2]
		if remoteName == "" || remoteMsg == "" {
			this.SendMsg("[im_server] command is not right, sendto:username:yourmassage exp: sendto|xxx|hello\n", this)
		}

		remoteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMsg("[im_server] user not find,please check!\n", this)
		} else {
			this.SendMsg("[private messages:"+this.Name+"] "+remoteMsg+"\n", remoteUser)
		}
	}
}

//用户的消息chan
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		_, err := this.conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("send massage err:" + this.Name)
			return
		}

	}

}

func (this *User) SendMsg(msg string, toUser *User) {
	toUser.C <- msg
	//toUser.conn.Write([]byte(msg))
}
