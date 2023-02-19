package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/rpc"
	"os"
	"time"
)

type Accepter struct {
	Port             string
	Final_string     string
	External_Storage string
}

func (a *Accepter) make_out(info string) {
	fmt.Println(info)
	a.Final_string = info
}
func (a *Accepter) check_if_over(l int) bool {
	f, err := os.Open(a.External_Storage)
	if err != nil {
		fmt.Println("外存打开错误:", err)
	}
	var number int = 0
	reader := bufio.NewReader(f)
	for {
		//多次循环读取
		_, err := reader.ReadString('\n') //以回车为分割依据将字符串读取
		if err == io.EOF {
			break
		}
		number++
	}
	if number == l {
		return true
	} //读取完毕后，如果外存数量达到要求行数就返回真
	return false
}
func (a *Accepter) check_if_empty(send string) bool {
	if send == "" {
		return true
	}
	return false
}
func (a *Accepter) Port_setter(t *bool, send *string) error {
	*t = true
	tcpAddr, err := net.ResolveTCPAddr("tcp", *send)
	if err != nil {
		fmt.Println("ResolveTCPAddr err=", err)
	}
	a.Port = tcpAddr.String()
	*send = "finished"
	return nil
}
func main() {
	a := Accepter{":8088", "dusha", "./out.TXT"} //默认端口8088
	//注册服务
	rpc.Register(a)
	label := false
	port := ":8088"
	a.Port_setter(&label, &port)
	//开始监听
	listener, err := net.Listen("tcp", a.Port)
	if err != nil {
		fmt.Println("监听错误！", err)
	}
	defer listener.Close()
	conn, err1 := listener.Accept()
	if err1 != nil {
		fmt.Println("Accept error:", err1)
	}
	defer conn.Close()
	for {
		//定时更新
		time.Sleep(5 * time.Second)
		buf := make([]byte, 1024)
		n, err2 := conn.Read(buf) //读取对方发送的信息
		if err != nil {
			if err2 == io.EOF {
				fmt.Println("决策信息接受完毕!")
			} else {
				log.Panicln("conn.Read err", err2)
				return
			}
		}
		a.make_out(string(buf[:n]))
		if a.check_if_empty(string(buf[:n])) {
			continue
		} //为空就不进行重置和写入
		//创建文件
		f, err := os.Create(a.External_Storage)
		if err != nil {
			log.Panicln(err)
		}
		//写入文件
		_, err = f.WriteString(string(buf))
		if err != nil && err != io.EOF {
			log.Panicln(err)
		}
		f.Close()
		time.Sleep(2 * time.Second)
		rpc.ServeConn(conn) //启用rpc服务
	}
}
