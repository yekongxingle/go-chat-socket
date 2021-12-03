package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

var ConnSlice map[net.Conn]*Heartbeat    //用来存储每一个连接

type Heartbeat struct {

	endTime int64
}

func main()  {
	ConnSlice=map[net.Conn]*Heartbeat{}   //初始化
	listen, err := net.Listen("tcp", "127.0.0.1:8080")
	if err!=nil{
		log.Fatal("服务器启动失败")
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err!=nil{
			fmt.Println("error accept:",err)
		}
		fmt.Printf("%s加入服务器\n",conn.RemoteAddr().String())
		ConnSlice[conn]=&Heartbeat{endTime: time.Now().Add(time.Second*5).Unix()}
		go handleConn(conn)
	}
}
func handleConn(conn net.Conn)  {
	//创建一个空间，用来存储连接中的数据
	buffer:=make([]byte,1024)
	//使用for循环，不停的对连接进行处理
	for {
		n, err := conn.Read(buffer)   //读连接中的数据
		if err!=nil{
			log.Fatal("error read")
		}
		//判断连接是否过期
		if ConnSlice[conn].endTime>time.Now().Unix(){
			//更新心跳时间
			ConnSlice[conn].endTime=time.Now().Add(time.Second*5).Unix()
		}else{
			log.Fatal("长时间未发消息断开连接")
		}
		//判断是否是心跳检测
		if string(buffer[:n])=="1"{    //如果是心跳检测，则不执行接下来的语句
			conn.Write([]byte("1"))
			continue
		}
		//对不是这次心跳检测的连接进行处理,如果过期就删除;!!!同时把这个连接的数据发给其他连接，这就是下面第一个判断的用处！！！
		for c,heart:=range ConnSlice{
			if c==conn{
				continue
			}
			if heart.endTime<time.Now().Unix(){
				delete(ConnSlice,c)
				c.Close()
				fmt.Println("delete connection ",c.RemoteAddr().String())
				fmt.Println("have connection :",ConnSlice)
				continue
			}
			//把数据发给其他连接
			c.Write(buffer[:n])
		}

	}
}