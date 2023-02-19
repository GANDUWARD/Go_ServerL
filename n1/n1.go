package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type data struct {
	Message  string
	Exponent string
	Timing   string
}

func print(w http.ResponseWriter, r *http.Request) {
	d := new(data)
	d.Message = "访问成功"
	d.Exponent = "拥挤"
	d.Timing = "80小时"
	d_json, _ := json.Marshal(d)
	s := string(d_json)
	fmt.Fprintf(w, s)
}
func main() {
	http.HandleFunc("/hello", print)
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}
