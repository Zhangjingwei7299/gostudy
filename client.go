package main

import (
	"fmt"
	"net"
	"flag"
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

var serverIp string
var serverPort int
// ./client -ip 127.0.0.1 -ip 8888
// 第一个是变量名，第二个是ip提示符，第三个是默认值
func init(){
	flag.StringVar(&serverIp,"ip","127.0.0.1","设置服务器IP地址(默认地址:127.0.0.1)")
	flag.IntVar(&serverPort,"port",8888,"设置服务器端口(默认:8888)")
}

func main(){
	//命令行解析
	flag.Parse()
	client:=NewClient(serverIp,serverPort)
	if client ==nil{
		fmt.Println(">>>>>>>>连接失败....")
		return
	}
	fmt.Println(">>>>连接服务器成功...")

	//启动客户端的业务
	select{}
}