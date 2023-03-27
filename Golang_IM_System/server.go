package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	// 在线用户的列表map
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//消息广播的channel
	Message chan string
}

// 创建一个server 接口(构造函数)
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// 监听Message广播消息channel中的goroutine，一旦有消息就发送给全部的在线User
func (this *Server) ListenMessager() {
	for {
		// 监听消息channel
		msg := <-this.Message

		// 1.将msg发送给全部的在线User
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// 广播消息的方法,两个参数,谁发起的,发送了什么消息
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "{" + user.Addr + "}" + user.Name + ":" + msg
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	// 当前链接业务 具体要执行的方法
	fmt.Println("链接建立成功")

	user := NewUser(conn, this)

	user.Online()

	// 监听用户是否活跃的channel
	isLive := make(chan bool)

	// 4. 接收客户端发送的消息
	go func() {
		buff := make([]byte, 4096)
		for {
			n, err := conn.Read(buff)
			if n == 0 { // 用户下线
				fmt.Println("{" + user.Addr + "}" + user.Name + " : Offline")
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
			}
			// 提取用户的消息（去除'\n'）
			msg := string(buff[:n-1])
			// 服务器日志
			fmt.Println("Receive message:" + msg + " from:{" + user.Addr + "}" + user.Name)

			// 用户针对msg进行消息处理
			user.DoMessage(msg)

			// 用户的任意消息 都代表当前用户是一个活跃用户
			isLive <- true
		}
	}()

	// 3.当前handler阻塞(防止主go程死亡)
	// 这里重置是因为select接收到一个管道后就结束了。外层的for使他一遍遍执行
	for {
		select {
		case <-isLive:
			// 说明当前用户是活跃的，应该重置定时器
			// 不做任何事情，为了激活select， 更新下面的定时器
			// 触发这个case后，后面的case都会重新执行
		case <-time.After(time.Second * 1000):
			// 说明已经超时，将当前的User强制关闭
			user.SendMsg("timeout,you have been offlined")
			//  销毁用的资源
			this.mapLock.Lock()
			delete(this.OnlineMap, user.Name)
			this.mapLock.Unlock()

			close(user.C)
			//// 关闭链接
			conn.Close()
			// 退出当前handler
			return // runtime.Goexit()
		}
	}
}

// 启动服务器端接口
func (this *Server) Start() {
	// 1.socket listen（以tcp形式监听ip和端口号）
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	// 4.close listen socket
	defer listener.Close() // 防止出现err后没有关闭listener

	// 启动监听Message的goroutine
	go this.ListenMessager()

	for {
		// 2.accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Listener accept err:", err)
			continue
		}
		// 3.do handler （开启一个协程处理请求，不影响其他连接请求）
		go this.Handler(conn)
	}
}
