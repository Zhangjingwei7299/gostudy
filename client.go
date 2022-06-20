package main

import (
	"fmt"
	"net"
)
type Client struct{
	ServerIp string
	ServerPort int
	Name string
	conn net.Conn

}

func NewClient(serverIp string,serverPort int)*Client{
	//创建客户端对象爱嗯
	client:=&Client{
		ServerIp:serverIp,
		ServerPort:serverPort,
	}

	//创建server
	conn,err:=net.Dial("tcp",fmt.Sprintf("%s:%d",serverIp,serverPort))
	if err!=nil{
		fmt.Println("net.Dial error:",err)
		return nil
	}

	client.conn=conn

	//返回对象爱
	return client
}

func main(){
	client:=NewClient("127.0.0.1",8888)
	if client ==nil{
		fmt.Println(">>>>>>>>连接失败....")
		return
	}
	fmt.Println(">>>>连接服务器成功...")

	//启动客户端的业务
	select{}
}