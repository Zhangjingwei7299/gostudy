package main

import "net"
type User struct{
	Name string
	Addr string
	C chan string
	conn net.Conn

	server *Server
}

//创建一个用户的API
func NewUser(conn net.Conn,server *Server)*User{
	userAddr :=conn.RemoteAddr().String()
	user :=&User{
		Name : userAddr,
		Addr : userAddr,
		C : make(chan string),
		conn : conn,
		server : server,
	}
	//一旦创建，就启动当前user channel消息的goroutione
	go user.ListenMessage()
	return user
}


func (this *User)Online(){
	//当前用户上线，广播，将用户先加入online表中

	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name]=this
	this.server.mapLock.Unlock()
	
	//广播当前用户上线信息
	this.server.BroadCast(this,"已上线")
}

func (this *User)Offline(){
	//当前用户下线，广播，将用户从online表中删除

	this.server.mapLock.Lock()
	// this.server.OnlineMap[this.Name]=this
	delete(this.server.OnlineMap,this.Name)
	this.server.mapLock.Unlock()
	
	//广播当前用户上线信息
	this.server.BroadCast(this,"下线")
}

//给当前user对应的客户端发送消息
func (this *User)SendMsg(msg string){
	this.conn.Write([]byte(msg))
}

//用户据处理消息的业务
func (this *User)DoMessage(msg string){
	if msg == "who"{
		//查询当前在线用户都有哪些

		this.server.mapLock.Lock()
		for _,user :=range this.server.OnlineMap{
			onlineMsg :="["+user.Addr+"]"+user.Name+":"+"在线...\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	}else{
		this.server.BroadCast(this,msg)
	}
	
}	




//监听当前User channel的方法,一旦有消息，就直接发送给对段客户端
func (this *User)ListenMessage(){
	for{
		msg:=<-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}