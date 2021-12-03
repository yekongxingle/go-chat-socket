package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main()  {
	//tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8080")
	//if err!=nil{
	//	log.Fatal("失败")
	//}
	//conn, err := net.DialTCP("tcp", nil, tcpAddr)
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err!=nil{
		log.Fatal(err.Error())
	}
	fmt.Println("connect success!")
	Sender(conn)
	fmt.Println("end")
}
func Sender(conn net.Conn)  {
	defer conn.Close()
	//创建一个stdin的reader
	reader := bufio.NewReader(os.Stdin)
	//创建一个发送心跳包的协程
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C  //每隔一秒，发送一个数据，不然会阻塞在这里，起到计时器的作用
			_, err := conn.Write([]byte("1"))
			if err!=nil{
				log.Fatal(err.Error())
			}
		}
	}()
	name:=""
	fmt.Print("请输入昵称：")
	fmt.Fscan(os.Stdin,&name)
	msg:=""
	buffer:=make([]byte,1024)
	timer := time.NewTimer(time.Second * 5) //创建一个定时器，每次服务器端发送数据就刷新时间
	defer timer.Stop()
	//创建一个判断服务器是否正常的协程
	go func() {
		<-timer.C
		log.Fatal("服务器出现故障，断开连接")
	}()
	//一个对计时器进行更新的协程
	go func() {
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				return
			}
			timer.Reset(time.Second * 5) //收到消息就刷新_t定时器，如果time.Second*5时间到了，那么就会<-_t.C就不会阻塞，代码会往下走，return结束
			if string(buffer[0:1]) != "1" { //心跳包消息定义为字符串"1",不需要打印出来
				fmt.Println(string(buffer[0:n]))
			}
		}
	}()
	//干正事了
	for {

		fmt.Fscan(reader, &msg)
		i := time.Now().Format("2006-01-02 15:04:05")
		conn.Write([]byte(fmt.Sprintf("%s--->%s: %s", i, name, msg))) //发送消息
	}
}