package ws

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestWsServer(t *testing.T) {

	// 启动websocket服务
	http.HandleFunc("/ws", wsHandle)
	err := http.ListenAndServe("9090", nil)
	if err != nil {
		log.Fatal("socket server run fail!", err)
	}

	select {}
}

func wsHandle(w http.ResponseWriter, r *http.Request) {
	var (
		resp []byte
		err  error
	)

	// 初始化新链接
	conn := NewWsServer(w, r)

	// 读取来信&业务分发
	for {
		if resp, err = conn.ReadMessage(); err != nil {
			goto ERR
		}
		fmt.Println(string(resp))

		// 可根据消息提的某个参数辨别具体业务，eg:resp.op=login
	}
ERR:
	conn.Close()
}
