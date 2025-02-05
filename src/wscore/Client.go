package wscore

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type WsShellClient struct {
	client *websocket.Conn
}

func NewWsShellClient(client *websocket.Conn) *WsShellClient {
	return &WsShellClient{client: client}
}

func (this *WsShellClient) Write(p []byte) (n int, err error) {
	fmt.Println("write...", string(p))
	err = this.client.WriteMessage(websocket.BinaryMessage,
		p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (this *WsShellClient) Read(p []byte) (n int, err error) {
	fmt.Println("read...p===", string(p))
	_, b, err := this.client.ReadMessage()
	fmt.Println("read...b===", string(b))

	if err != nil {
		return 0, err
	}
	return copy(p, string(b)), nil
}
