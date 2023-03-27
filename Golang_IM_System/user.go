package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	// 表明用户属于哪一个Server
	server *Server
}

// 创建一个用户的API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	// 启动监听当前User channel消息的goroutine
	go user.ListenMessage()

	return user
}

// 用户上线业务
func (this *User) Online() {
	// 1.用户上线，将用户加入到Server的OnlineMap中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	// 2.广播当前用户上线消息
	this.server.BroadCast(this, "Online")
}

// 用户下线业务
func (this *User) Offline() {
	// 1.用户下线
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	// 2.广播当前用户上线消息
	this.server.BroadCast(this, "Offline")
}

// 给当前User用户对应的客户端发送消息
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

// 用户处理消息的业务
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		// ”who“ 查询当前在线用户都有哪些
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "Online\n"
			this.SendMsg(onlineMsg)
			// 或者直接：user.C <- onlineMsg
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// rename| 张三 ：  修改用户名
		newname := strings.Split(msg, "|")[1]
		// 判断name是否存在
		if _, ok := this.server.OnlineMap[newname]; ok {
			// 查询成功，说明已经有人取这个名了，不能取
			this.SendMsg("this name has been used\n")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newname] = this
			this.server.mapLock.Unlock()

			this.Name = newname
			this.SendMsg("Update username sucessful : " + newname + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// 消息格式：to|张三|消息内容
		// 1. 获取对方的用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.SendMsg("message format is invalid\n")
			return
		}
		// 2. 根据用户名,得到对方User对象
		remoteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMsg("User:" + remoteName + "is not existed\n")
			return
		}
		// 3. 获取消息内容，通过对方的User对象将消息内容发送过去
		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.SendMsg("content must not be blanked\n")
			return
		}
		remoteUser.SendMsg(this.Name + " send to you : " + content)
	} else {
		// 将得到的消息进行广播
		this.server.BroadCast(this, msg)
	}
}

// 监听当前User channel的方法，一旦有消息，就直接发送给对应客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		this.conn.Write([]byte(msg + "\n"))
	}
}
