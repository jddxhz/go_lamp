package main

import (
	"fmt"
	"net"
	"net/http"

	sendMail "code.lampworld.xyz/go_lamp/websocket/mail"
	"code.lampworld.xyz/go_lamp/websocket/myweb"
	"code.lampworld.xyz/go_lamp/websocket/socket"

	"code.lampworld.xyz/go_lamp/websocket/timer"
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
	getHostAddr()

	//将/js/路径下的请求匹配到 ./static/js/下
	http.Handle("/img/", http.FileServer(http.Dir("./view")))
	http.Handle("/css/", http.FileServer(http.Dir("./view")))
	http.Handle("/js/", http.FileServer(http.Dir("./view")))

	fmt.Println("服务器即将开启，访问地址 http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("服务器开启错误: ", err)
	}
}

func getHostAddr() {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						sendMail.SendMailByNetSMTP("神灯赏金任务运行于 " + ipnet.IP.String())
					}
				}
			}
		}
	}
}
