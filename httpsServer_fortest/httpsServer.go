package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type MyServer struct {
	SServer      http.Server
	Final_String string
}

func (s *MyServer) hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(s.Final_String + "0.125;十五分钟"))
}
func (s *MyServer) socket_send() {
	//用于接受预测模型通过socket发送的数据，并对其双向检验
	l, err := net.Listen("tcp", ":44329")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	conn, err := l.Accept()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("接受成功！")
	}
	defer conn.Close()
	var buf []byte
	conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("接受成功！")
	}
	s.Final_String = string(buf)
	fmt.Println(s.Final_String)
	_, err = conn.Write([]byte(s.Final_String))
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("检验发送内容成功！")
	}
}
func (s *MyServer) WriteToFront() {
	//对前端发FinalString
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(s.Final_String))
		})
		if err := s.SServer.ListenAndServeTLS("ganduward.com_bundle.crt", "ganduward.com.key"); err != http.ErrServerClosed {
			log.Panicln("前端发送错误!", err)
		}
	}()
}
func (s *MyServer) Close() {
	err := s.SServer.Shutdown(nil)
	if err != nil {
		log.Panic("关闭错误！", err)
	}
}
func main() {
	TheServer := MyServer{http.Server{Addr: ":443"}, ""}
	TheServer.WriteToFront()
	for {
		//TheServer.socket_send()
		time.Sleep(5 * time.Second)
	}
}
