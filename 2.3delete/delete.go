package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	url := "https://www.4399.com"
	payload := strings.NewReader("{\"userId\":1,\"articleId\":1,\"comment\":\"这是一条评论\"}")
	req, _ := http.NewRequest("DELETE", url, payload)
	req.Header.Add("Content-Type", "application/json")
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(res)
	fmt.Println(string(body))
}
