package socket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"code.lampworld.xyz/go_lamp/websocket/impl"
	"code.lampworld.xyz/go_lamp/websocket/task"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func WsHandler(w http.ResponseWriter, r *http.Request) {
	type login struct {
		Account string  `json:"account"`
		Cipher  string  `json:"cipher"`
		Sdc     float64 `json:"sdc"`
		TeamID  int     `json:"team"`
	}
	//	w.Write([]byte("hello"))
	var (
		wsConn *websocket.Conn
		err    error
		conn   *impl.Connection
		data   []byte
		l      login
		pass   string
		ip     string
	)
	if err != nil {
		fmt.Println(err)
	}
	tort := task.GetTort()
	// 完成ws协议的握手操作
	// Upgrade:websocket
	if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil {
		return
	}
	ip = fmt.Sprintf("%v", wsConn.RemoteAddr())
	fmt.Printf("%v 连接成功...\n", ip)

	if conn, err = impl.InitConnection(wsConn); err != nil {
		goto ERR
	}

	// 启动线程，不断发消息
	go func() {
		for {
			tortInfo := tort.TortInfo()
			if err != nil {
				return
			}
			if err = conn.WriteMessage([]byte(tortInfo)); err != nil {
				return
			}
			if tort.IsGameStart {
				time.Sleep(40 * time.Millisecond)
			} else {
				time.Sleep(1 * time.Second)
			}
		}
	}()

	for {
		data, err = conn.ReadMessage()
		if len(data) < 6 {
			goto ERR
		}
		t := string(data)[:6]
		// fmt.Println(string(data[6:]))
		switch t {
		case "login-":
			err = json.Unmarshal(data[6:], &l)
			p, err := task.Login(l.Account, l.Cipher)
			if err != nil {
				pass = err.Error()
			} else {
				pass = fmt.Sprintf(`{"LoginResult":{"code":0,"msg":"成功","vip":%v,"sdc":%v,"team":%v}}`, p.Vip, p.Sdc, p.TeamID)
			}
		case "jion--":
			err = json.Unmarshal(data[6:], &l)
			p, err := task.JionTeam(l.Account, l.Cipher, l.Sdc, l.TeamID)
			if err != nil {
				pass = err.Error()
			} else {
				pass = fmt.Sprintf(`{"LoginResult":{"code":0,"msg":"成功","vip":%v,"sdc":%v,"team":%v}}`, p.Vip, p.Sdc, p.TeamID)
			}
		case "award-":
			err = json.Unmarshal(data[6:], &l)
			p, _ := task.Login(l.Account, l.Cipher)
			err := p.Award()
			pass = err.Error()
		case "000001":
			fmt.Println(string(data[:6]))
		}
		if err != nil {
			goto ERR
		}
		if err = conn.WriteMessage([]byte(pass)); err != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()
	fmt.Printf("%v 断开连接...\n", ip)
}
