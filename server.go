package main

import (
	"net"
	"fmt"
	"sync"
	"io"
)
type Server struct{
	Ip string
	Port int

	//在线用户的列表,属于公共资源，要注意枷锁
	OnlineMap map[string]*User
	//锁
	mapLock sync.RWMutex

	//消息广播的channel
	Message chan string
	
}

//创建一个server的接口
func NewServer(ip string,port int)*Server{
	server:=&Server{
		Ip:ip,
		Port:port,
		OnlineMap:make(map[string]*User),
		Message : make(chan string),
	}
	return server
}


//监听Message广播消息channel的goroutine，一旦有消息就发送给全部在线User
func (this *Server)ListenMessager(){
	for{
		msg:=<-this.Message

		//将msg发送给全部在线User
		this.mapLock.Lock()
		for _,cli:=range this.OnlineMap{
			cli.C<-msg
		}
		this.mapLock.Unlock()
	}
}


//广播消息的方法
func (this *Server)BroadCast(user *User,msg string){
	sendMsg:="["+user.Addr+"]"+user.Name+":"+msg

	this.Message<-sendMsg
}

func (this *Server)Handler(conn net.Conn){
	//...当前连接的业务
	//fmt.Println("连接建立成功")
	user:=NewUser(conn)
	//当前用户上线，广播，将用户先加入online表中
	this.mapLock.Lock()
	this.OnlineMap[user.Name]=user
	this.mapLock.Unlock()

	//广播当前用户上线信息
	this.BroadCast(user,"已上线")

	//接受客户端发送的消息
	go func(){
		buf :=make([]byte,4096)
		for{
			n,err := conn.Read(buf)
			if n==0{
				this.BroadCast(user,"下线")
				return
			}
			if err!=nil && err!=io.EOF{
				fmt.Println("Conn Read err:",err)
				return
			}
			
			//提取用户消息去掉\n
			msg := string(buf[:n-1])
			
			//将得到的消息进行广播
			this.BroadCast(user,msg)
		}
	}()
	//当前handler阻塞
	select {}
}

//启动服务器的接口
func (this *Server)Start(){
	//scoket listen
	listener,err:=net.Listen("tcp",fmt.Sprintf("%s:%d",this.Ip,this.Port))
	if err!=nil{
		fmt.Println("net.Listen err:",err)
	}

	//close listen socket
	defer listener.Close()

	//启动监听Message的goroutine
	go this.ListenMessager()

	for {
		//accept
		conn,err:=listener.Accept()
		if err!=nil{
			fmt.Println("listener accept err:",err)
			continue
		}
		
		//do handler
		go this.Handler(conn)
	}
	


	
}