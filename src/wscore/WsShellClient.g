package wscore

import (
	"fmt"
	"github.com/gorilla/websocket"
	"strings"
)

type WsShellClient struct {
	client *websocket.Conn
}

func NewWsShellClient(client *websocket.Conn) *WsShellClient {
	return &WsShellClient{client: client}
}



func (this *WsShellClient) Read(p []byte) (n int, err error) {
	// 读取前端传来的数据
	_, b, err := this.client.ReadMessage()
	if err != nil {
		return 0, err
	}

	// 过滤掉转义序列（类似 \x1b[2;5R 这样的内容）
	command := string(b)
	if strings.HasPrefix(command, "\x1b") {
		return 0, nil // 忽略转义序列
	}

	// 确保命令以换行符结束
	if len(command) > 0 && !strings.HasSuffix(command, "\n") {
		command += "\n"
	}

	fmt.Printf("Processed command: %q\n", command)
	return copy(p, command), nil
}

func (this *WsShellClient) Write(p []byte) (n int, err error) {
	// 过滤掉可能的命令回显
	output := string(p)
	fmt.Println("write to client: ", output)
	// 发送到客户端
	err = this.client.WriteMessage(websocket.TextMessage, []byte(output))
	if err != nil {
		return 0, err
	}
	return len(p), nil
}


