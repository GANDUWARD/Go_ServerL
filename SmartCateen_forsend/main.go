package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type MyServer struct {
	SServer          http.Server
	Final_strings    [3][3]string //用于存储各食堂各层的结果
	External_Storage string
	router           *gin.Engine
}

func (s *MyServer) check_if_over(l int) bool {
	f, err := os.Open(s.External_Storage)
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
func (s *MyServer) Manage_strings() {
	if !s.check_if_over(9) {
		fmt.Println("数据未收集完毕！")
		return
	}
	f := s.Open_External_Storage()
	defer f.Close()
	reader := bufio.NewReader(f)
	for {
		//多次循环读取
		line, err := reader.ReadString('\n')   //以回车为分割依据将字符串读取
		real_line := strings.Split(line, "\n") //将回车分割出来
		if err == io.EOF {
			break
		}
		var m int = int(real_line[0][0]) - 49
		var n int = int(real_line[0][2]) - 49
		s.Final_strings[m][n] = real_line[0][4:] + ";"

	}

}
func (s *MyServer) Open_External_Storage() *os.File {
	f, err := os.OpenFile(s.External_Storage, os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	return f
}
func (s *MyServer) Serve_justby_pkg_http() {
	http.HandleFunc("/cateen1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			s.Manage_strings()
			w.Write([]byte(s.Final_strings[0][0] + s.Final_strings[0][1] + s.Final_strings[0][2]))
		} else if r.Method == "POST" {
			fmt.Println("Click!")
		}
	})
	http.HandleFunc("/cateen1/cateen2", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			s.Manage_strings()
			w.Write([]byte(s.Final_strings[1][0] + s.Final_strings[1][1] + s.Final_strings[1][2]))
		} else if r.Method == "POST" {
			fmt.Println("Click!")
		}
	})
	http.HandleFunc("/cateen1/cateen2/cateen3", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			s.Manage_strings()
			w.Write([]byte(s.Final_strings[2][0] + s.Final_strings[2][1] + s.Final_strings[2][2]))
		} else if r.Method == "POST" {
			fmt.Println("Click!")
		}
	})
	err := s.SServer.ListenAndServeTLS("ganduward.com_bundle.crt", "ganduward.com.key")
	if err != http.ErrServerClosed {
		panic(err) // 为了捕获监听失败的错误
	}
}
func (s *MyServer) cateen1(c *gin.Context) {
	s.Manage_strings()
	c.String(http.StatusOK, s.Final_strings[0][0]+s.Final_strings[0][1]+s.Final_strings[0][2])
}
func (s *MyServer) cateen2(c *gin.Context) {
	s.Manage_strings()
	c.String(http.StatusOK, s.Final_strings[1][0]+s.Final_strings[1][1]+s.Final_strings[1][2])
}
func (s *MyServer) cateen3(c *gin.Context) {
	s.Manage_strings()
	c.String(http.StatusOK, s.Final_strings[2][0]+s.Final_strings[2][1]+s.Final_strings[2][2])
}
func (s *MyServer) WebToFront() {
	s.Manage_strings()
	s.router = gin.Default()
	v2 := s.router.Group("/")
	{
		v2.GET("/cateen1", s.cateen1)
		v2.GET("/cateen2", s.cateen2)
		v2.GET("/cateen3", s.cateen3)

	}
	s.router.Run(":8088")
	//s.Serve_justby_pkg_http()
}
func (s *MyServer) Close() error {
	err := s.SServer.Shutdown(nil)
	return err
}
func main() {
	ser := MyServer{http.Server{Addr: ":443"}, [3][3]string{}, "/root/Go_works/SmartCateen_foraccept/out.TXT", &gin.Engine{}}
	ser.WebToFront()
}
