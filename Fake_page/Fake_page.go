package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

func greet(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello World! %s", time.Now())
	t, err := template.ParseFiles("./template_example.tmpl")
	if err != nil {
		fmt.Println("template parsefile failed,err:", err)
		return
	}
	// 2.渲染模板
	name := "dusha"
	t.Execute(w, name)
}
func show_picture(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("./image/image.jpg")
	if err != nil {
		fmt.Println("文件打开错误！", err)
	}
	var b []byte
	f.Read(b)
	f.Close()
	w.Write(b)
}
func main() {
	//http.Handle("/image", http.FileServer(http.Dir("./image/image.jpg"))) //将图片“服务”上去
	http.HandleFunc("/", greet)
	//http.HandleFunc("/image", show_picture)
	http.HandleFunc("/image/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	}) //注意这里要额外提供一个空间把你前面tmpl填充文件通过这个ServeFile读取指定路径并将文件服务上去
	//实现高并发编程
	for {
		err := http.ListenAndServeTLS(":443", "ganduward.com_bundle.crt", "ganduward.com.key", nil)
		if err != nil {
			fmt.Println(err)
		}
	}
}
