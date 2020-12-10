package ws

import (
	"com/utils"
	"fmt"
	"testing"
)

func TestClient(t *testing.T) {
	conn, err := NewWsClient("wss://api.hbdm.com/linear-swap-ws")

	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		data, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			continue
		}

		res, err := utils.GZipDecompress(data)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("消息：", string(res))
	}
}
