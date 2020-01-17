package main

import (
	"fmt"
	"net/http"

	"go_code/websocket/myweb"
	"go_code/websocket/socket"

	"go_code/websocket/timer"
)

func main() {
	go timer.Timer()

	//ws 协议的 websocket 服务
	go func() {
		http.HandleFunc("/ws", socket.WsHandler)
		err := http.ListenAndServe("0.0.0.0:7777", nil)
		if err != nil {
			fmt.Println("服务器开启错误: ", err)
		}
	}()

	//http 协议的 网页服务
	http.HandleFunc("/", myweb.MyWeb)

	//将/js/路径下的请求匹配到 ./static/js/下
	http.Handle("/img/", http.FileServer(http.Dir("./view")))
	http.Handle("/css/", http.FileServer(http.Dir("./view")))
	http.Handle("/js/", http.FileServer(http.Dir("./view")))

	fmt.Println("服务器即将开启，访问地址 http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("服务器开启错误: ", err)
	}
}
