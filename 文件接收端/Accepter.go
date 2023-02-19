package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	//监听
	listener, err := net.Listen("tcp", ":8088")
	if err != nil {
		fmt.Println("net.Listen err", err)
		return
	}

	defer listener.Close()

	//堵塞等待用户连接
	conn, err1 := listener.Accept()
	if err1 != nil {
		fmt.Println("listener.Accept err:", err1)
		return
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	n, err2 := conn.Read(buf) //读取对方发送的文件名
	if err2 != nil {
		fmt.Println("conn.Read err", err2)
		return
	}
	fileName := string(buf[:n])
	//回复“ok”
	conn.Write([]byte("ok"))

	//接受文件内容
	RecvFile(fileName, conn)
}

func RecvFile(fileName string, conn net.Conn) {
	//新建文件
	f, err3 := os.Create(fileName)
	if err3 != nil {
		fmt.Println("os.Create err:", err3)
		return
	}

	buf := make([]byte, 1024*500)
	//接受多少，写多少
	for {
		n, err4 := conn.Read(buf) //接受对方发送过来的文件内容
		if err4 != nil {
			if err4 == io.EOF {
				fmt.Println("文件接受完毕")
			} else {
				fmt.Println("conn.Read err", err4)
			}
			return
		}
		if n == 0 {
			fmt.Println("n == 0,文件接受完毕")
			return
		}
		f.Write(buf[0:n])
	}
}
