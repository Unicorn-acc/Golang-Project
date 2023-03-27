> 学习时间：23.3.23
>
> 视频：[https://www.bilibili.com/video/BV1gf4y1r79E](https://www.bilibili.com/video/BV1gf4y1r79E)

# 【Golang 快速入门4】项目实战：即时通信系统

# 即时通信系统 - 服务端

项目架构图：

![](https://img-blog.csdnimg.cn/img_convert/2e8d017c02043293d8b4d60899afc8a4.png)

9个小版本迭代：

- 版本一：构建基础 Server
- 版本二：用户上线功能
- 版本三：用户消息广播机制
- 版本四：用户业务层封装
- 版本五：在线用户查询
- 版本六：修改用户名
- 版本七：超时强踢功能
- 版本八：私聊功能
- 版本九：客户端实现

## 版本一：构建基础 Server

 server.go，其中包含以下内容：

- 定义 Server 结构体，包含 IP、Port 字段
- `NewServer(ip string, port int)` 创建 Server 对象的方法
- `(s *Server) Start()` 启动 Server 服务的方法
- `(s *Server) Handler(conn net.Conn)` 处理连接业务

```go
type Server struct {
	Ip   string
	Port int
}

// 创建一个server 接口(构造函数)
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

func (this *Server) Handler(conn net.Conn) {
	// 当前链接业务 具体要执行的方法
	fmt.Println("链接建立成功")
}

// 启动服务器端接口
func (this *Server) Start() {
	// socket listen（以tcp形式监听ip和端口号）
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	// close listen socket
	defer listener.Close() // 防止出现err后没有关闭listener
	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Listener accept err:", err)
			continue
		}
		// do handler （开启一个协程处理请求，不影响其他连接请求）
		go this.Handler(conn)
	}

}

```

main.go，启动我们编写的 Server：

```go
package main

func main() {
	// 都属于main包，不需要import
	server := NewServer("127.0.0.1", 8888)
	server.Start()
}

```

window下编译运行：

同时编译编写的两个文件：`go build -o server.exe`

然后运行编译出的文件：`server.exe`

使用命令侦听我们构建的服务：`nc 127.0.0.1 8888`

## 版本二：用户上线+广播功能

![](https://img-blog.csdnimg.cn/img_convert/3e91499d6a9d0543cec3d8523269bad7.png)

user.go：

- `NewUser(conn net.Conn) *User` 创建一个 user 对象
- `(u *User) ListenMessage()` 监听 user 对应的 channel 消息

```go
type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

// 创建一个用户的API
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}
	// 启动监听当前User channel消息的goroutine
	go user.ListenMessage()

	return user
}

// 监听当前User channel的方法，一旦有消息，就直接发送给对应客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		this.conn.Write([]byte(msg + "\n"))
	}
}
```

Server.go：

- 新增 OnlineMap 和 Message 属性
- 在处理客户端上线的 Handler 创建并添加用户
- 新增广播消息方法
- 新增监听广播消息 channel 方法
- 用一个 goroutine 单独监听 Message

```go
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

	user := NewUser(conn)

	// 1.用户上线，将用户加入到OnlineMap
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	// 2.广播当前用户上线消息
	this.BroadCast(user, "已上线")

	// 3.当前handler阻塞(防止主go程死亡)
	select {}
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
```

学习到的编程思路：

- 结构体中的 channel 基本都需要开个循环去监听其变化（尝试取出值，发送给其他 channel）

> 总结：就是每个用户上线的时候，就会新建一个用户的对象（结构体），对象中自带一个goroutine，然后每个用户上线的时候，都把这个消息遍历发送给每个用户的goroutine
>
>  一、user.go 后端服务器用来当前用户的类型
>
> ​	1.创建一个user对象
>
> 	2. 监听user对应的channel消息
>
> 二、server.go 
>
> 	1. 结构体新增OnlineMap和Message属性
> 	1. 在处理客户端上线的Handler创建并添加用户，并广播消息(消息送入chan)
> 	1. 新增广播消息的方法
> 	1. 新增监听广播消息channel的方法，得到消息发给OnlineMap中每一个User的chan中
> 	1. 在Start方法中用一个goroutine单独监听Message

## 版本三：用户消息广播机制

server.go：完善 handle 处理业务方法，启动一个针对当前客户端的读 routine

![](C:\Users\Admin\Desktop\Golang\Golang_CommunicateSystem.assets\image-20230323141419535.png)

server.go

```go
package main

import (
	"fmt"
	"io"
	"net"
	"sync"
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

	user := NewUser(conn)

	// 1.用户上线，将用户加入到OnlineMap
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	// 2.广播当前用户上线消息
	this.BroadCast(user, "Online")

	// 4. 接收客户端发送的消息
	go func() {
		buff := make([]byte, 4096)
		for {
			n, err := conn.Read(buff)
			if n == 0 { // 用户下线
				this.BroadCast(user, "Offline")
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
			}
			// 提取用户的消息（去除'\n'）
			msg := string(buff[:n-1])
			// 服务器日志
			fmt.Println("Receive message:" + msg + " from:{" + user.Addr + "}" + user.Name)
			// 将得到的消息进行广播
			this.BroadCast(user, msg)
		}
	}()

	// 3.当前handler阻塞(防止主go程死亡)
	select {}
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
```

## 版本四：用户业务层封装

用户上线、用户下线、用户处理消息都是User的业务，不应该是Server的业务

user.go：

- user 类型新增 server 关联
- 新增 Online、Offline、DoMessage 方法

```go
package main

import "net"

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

// 用户处理消息的业务
func (this *User) DoMessage(msg string) {
	// 将得到的消息进行广播
	this.server.BroadCast(this, msg)
}

// 监听当前User channel的方法，一旦有消息，就直接发送给对应客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		this.conn.Write([]byte(msg + "\n"))
	}
}
```

server.go：

- 使用 user 封装好的业务替换之前的代码

```go
package main

import (
	"fmt"
	"io"
	"net"
	"sync"
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

	// 4. 接收客户端发送的消息
	go func() {
		buff := make([]byte, 4096)
		for {
			n, err := conn.Read(buff)
			if n == 0 { // 用户下线
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
		}
	}()

	// 3.当前handler阻塞(防止主go程死亡)
	select {}
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
```

> 总结：每个模块处理各自的事情，因此将用户上线、用户下线、用户发送消息封装到user.go中去
>
> user.go：结构体新增server关联
>
> 	1. 新增Online方法
> 	1. 新增Offline方法
> 	1. 新增DoMessage方法
>
> server.go：将之前user的业务进行替换

## 版本五：在线用户查询

若某个用户输入的消息为 `who` 则查询当前在线用户列表。

user.go：

- 提供 SendMsg 向对象客户端发送消息 API

```go
// 给当前User用户对应的客户端发送消息
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

```

- 在 DoMessage() 方法中，加上对 “who” 指令的处理，返回在线用户信息

```go
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
	} else {
		// 将得到的消息进行广播
		this.server.BroadCast(this, msg)
	}
}
```

## 版本六：修改用户名

若某个用户输入的消息为 `rename|张三` 则将自己的 Name 修改为张三。

user.go：

- 在 DoMessage() 方法中，加上对 “rename|张三” 指令的处理

```go
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
	} else {
		// 将得到的消息进行广播
		this.server.BroadCast(this, msg)
	}
}
```

> 这里存在漏洞,rename||123,用户名会为空，应当将用户名中的|进行转义

## 版本七：超时强踢功能

用户的任意消息表示用户为活跃，长时间不发消息认为超时，就才一强制关闭用户连接。

server.go：

- 在用户 Handler() goroutine 中，添加活跃用户 channel，一旦用户有消息，就向该 channel 发送数据

```go
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
		case <-time.After(time.Second * 10):
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
```

## 版本八：私聊功能

消息格式：`to|张三|你好啊，我是...`

user.go，在 DoMessage() 方法中，加上对 “to|张三|你好啊” 指令的处理：

```go
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

```

## 版本九：客户端实现

# 即时通信系统 - 客户端

> 以下代码都是在 client.go 文件中

## 客户端类型定义与链接

client.go：

```go

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}
	// 链接Server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err : ", err)
		return nil
	}
	client.conn = conn

	// 返回对象
	return client
}

func main() {
	clint := NewClient("127.0.0.1", 8888)
	if clint == nil {
		fmt.Println(">>>>>> Connect Server failed")
		return
	}
	fmt.Println(">>>>>> Connect Server success")

	// 启动客户端的业务
	select {}
}
```

编译指令：`go build -o client.exe client.go`

运行编译后的文件：`./client`

## 解析命令行 flag

在 init 函数中初始化命令行参数并解析：

```go
// 绑定命令行参数
var serverIp string
var serverPort int

func init() {
	// ./client -ip 127.0.0.1 -port 8888
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认是127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口(默认是8888)")
}

func main() {
	// 命令行解析
	flag.Parse()
	//clint := NewClient("127.0.0.1", 8888)
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>> Connect Server failed")
		return
	}
	fmt.Println(">>>>>> Connect Server success")

	// 启动客户端的业务
	select {}
}
```

然后在运行客户端时可以通过 命令行传参运行：

```bash
./client -ip 127.0.0.1 -port 8888
```

![](C:\Users\Admin\Desktop\Golang\Golang_CommunicateSystem.assets\image-20230323214507402.png)

## 菜单显示

给 Client 新增 flag 属性：

```go
type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int // 当前客户端的模式
}
```

新增 menu() 方法，获取用户输入的模式：

```go
// 菜单
func (client *Client) menu() bool {
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>请输入合法范围内的数字<<<<")
		return false
	}
}
```

新增 Run() 主业务循环：

```go
func (client *Client) Run() {
	for client.flag != 0 { // 等待用户输入模式
		for client.menu() != true {
		}

		// 根据不同模式处理不同的业务
		switch client.flag {
		case 1:
			// 公聊模式
			fmt.Println("公聊模式")
		case 2:
			// 私聊模式
			fmt.Println("私聊模式")
		case 3:
			// 更新用户名
			fmt.Println("更新用户名")

		}
	}
	fmt.Println("退出")
}
```

client.go

```go
type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int // 当前客户端的模式
}

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	// 链接Server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err : ", err)
		return nil
	}
	client.conn = conn

	// 返回对象
	return client
}

// 菜单显示
func (client *Client) menu() bool {
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>请输入合法范围内的数字<<<<")
		return false
	}
}

func (client *Client) Run() {
	for client.flag != 0 { // 等待用户输入模式
		for client.menu() != true {
		}

		// 根据不同模式处理不同的业务
		switch client.flag {
		case 1:
			// 公聊模式
			fmt.Println("公聊模式")
		case 2:
			// 私聊模式
			fmt.Println("私聊模式")
		case 3:
			// 更新用户名
			fmt.Println("更新用户名")

		}
	}
	fmt.Println("退出")
}

// 绑定命令行参数
var serverIp string
var serverPort int

func init() {
	// ./client -ip 127.0.0.1 -port 8888
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认是127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口(默认是8888)")
}

func main() {
	// 命令行解析
	flag.Parse()
	//clint := NewClient("127.0.0.1", 8888)
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>> Connect Server failed")
		return
	}
	fmt.Println(">>>>>> Connect Server success")

	// 启动客户端的业务
	client.Run()
}

```

## 更新用户名

新增 UpdateName() 更新用户名：

```go
func (client *Client) UpdateName() bool {
	fmt.Println(">>> 请输入用户名:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n" // 封装协议
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err: ", err)
		return false
	}
	return true
}
```

添加 server 回执消息方法 DealResponse()

```go
// 处理server回应的消息，直接显示到标准输出即可
func (client *Client) DealResponse() {
	// 一旦client.conn有数据，直接copy到stdout标准输出上，永久阻塞监听
	io.Copy(os.Stdout, client.conn)
	// 上面这段等价于：
	//for {
	//	buf := make([]byte, 4096)
	//	client.conn.Read(buf)
	//	fmt.Println(buf)
	//}
}
```

在 main 中开启一个 goroutine，去承载 DealResponse() 流程：

```go
func main() {
	// 命令行解析
	flag.Parse()
	//clint := NewClient("127.0.0.1", 8888)
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>> Connect Server failed")
		return
	}
	fmt.Println(">>>>>> Connect Server success")

	// 单独开启一个goroutine处理server的回执消息
	go client.DealResponse()
	// 启动客户端的业务
	client.Run()
}

```

## 公聊模式

新增 PublicChat() 公聊模式：

```go
func (client *Client) PublicChat() {
	// 提示用户输入消息
	var chatMsg string

	fmt.Println(">>>>请输入聊天内容，exit退出.")
	fmt.Scanln(&chatMsg)

	// 发送给服务器
	for chatMsg != "exit" {
		// 发给服务器
		// 消息不为空立即发送
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err: ", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println(">>>>请输入聊天内容，exit退出.")
		fmt.Scanln(&chatMsg)
	}
}
```

## 私聊模式

查询当前有哪些用户在线：

```go
// 查询在线用户
func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return
	}
}
```

新增私聊业务：

```go
func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.SelectUsers()
	fmt.Println(">>>>请输入聊天对象的[用户名], exit退出: ")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>请输入消息内容，exit退出:")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			// 消息不为空则发送
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"

				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					return
				}
				chatMsg = ""
				fmt.Println(">>>>请输入消息内容，exit退出:")
				fmt.Scanln(&chatMsg)
			}
		}
		client.SelectUsers()
		fmt.Println(">>>>请输入聊天对象的[用户名], exit退出: ")
		fmt.Scanln(&remoteName)
	}
}
```

# 三个文件完整代码

## server.go

```go
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

```

## main.go

```go
package main

func main() {
	// 都属于main包，不需要import
	server := NewServer("127.0.0.1", 8888)
	server.Start()
}
```

## client.go

```go
package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int // 当前客户端的模式
}

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	// 链接Server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err : ", err)
		return nil
	}
	client.conn = conn

	// 返回对象
	return client
}

// 处理server回应的消息，直接显示到标准输出即可
func (client *Client) DealResponse() {
	// 一旦client.conn有数据，直接copy到stdout标准输出上，永久阻塞监听
	io.Copy(os.Stdout, client.conn)
	// 上面这段等价于：
	//for {
	//	buf := make([]byte, 4096)
	//	client.conn.Read(buf)
	//	fmt.Println(buf)
	//}
}

// 菜单显示
func (client *Client) menu() bool {
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>请输入合法范围内的数字<<<<")
		return false
	}
}

func (client *Client) UpdateName() bool {
	fmt.Println(">>> 请输入用户名:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n" // 封装协议
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err: ", err)
		return false
	}
	return true
}

func (client *Client) PublicChat() {
	// 提示用户输入消息
	var chatMsg string

	fmt.Println(">>>>请输入聊天内容，exit退出.")
	fmt.Scanln(&chatMsg)

	// 发送给服务器
	for chatMsg != "exit" {
		// 发给服务器
		// 消息不为空立即发送
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err: ", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println(">>>>请输入聊天内容，exit退出.")
		fmt.Scanln(&chatMsg)
	}
}

// 查询在线用户
func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.SelectUsers()
	fmt.Println(">>>>请输入聊天对象的[用户名], exit退出: ")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>请输入消息内容，exit退出:")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			// 消息不为空则发送
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"

				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					return
				}
				chatMsg = ""
				fmt.Println(">>>>请输入消息内容，exit退出:")
				fmt.Scanln(&chatMsg)
			}
		}
		client.SelectUsers()
		fmt.Println(">>>>请输入聊天对象的[用户名], exit退出: ")
		fmt.Scanln(&remoteName)
	}
}

func (client *Client) Run() {
	for client.flag != 0 { // 等待用户输入模式
		for client.menu() != true {
		}

		// 根据不同模式处理不同的业务
		switch client.flag {
		case 1:
			// 公聊模式
			client.PublicChat()
		case 2:
			// 私聊模式
			client.PrivateChat()
		case 3:
			// 更新用户名
			client.UpdateName()
		}
	}
	fmt.Println("退出")
}

// 绑定命令行参数
var serverIp string
var serverPort int

func init() {
	// ./client -ip 127.0.0.1 -port 8888
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认是127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口(默认是8888)")
}

func main() {
	// 命令行解析
	flag.Parse()
	//clint := NewClient("127.0.0.1", 8888)
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>> Connect Server failed")
		return
	}
	fmt.Println(">>>>>> Connect Server success")

	// 单独开启一个goroutine处理server的回执消息
	go client.DealResponse()
	// 启动客户端的业务
	client.Run()
}
```

