package web_socket

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

const (
	subscribeType = 1
)

type WebSocket interface {
	Subscribe() chan []byte
	Err() chan error
	Start(ctx context.Context)
	UnsubscribeMessage() string
}

type webSocket struct {
	ctx     context.Context   //线程管理组件
	url     string            //websocket地址
	method  string            //需要额外发送的信息
	resp    string            //用来取消订阅的消息体，需要用户自行分析
	dialer  *websocket.Dialer //拨号配置信息
	conn    *websocket.Conn   //websocket客户端
	channel chan []byte       //消息传递信道
	err     chan error        //错误信息传递信道
}

/**
** 连接对端websocket服务，url为服务地址，method为需要传送的方法，可以为空，表示不需要传递信息到对端开启服务
 */
func NewWebSocket(ctx context.Context, url, method string) (WebSocket, error) {
	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
		Subprotocols:     []string{"ws", "wss"},
	}
	cli, resp, err := dialer.DialContext(ctx, url, nil)
	if err != nil {
		return nil, errors.New(url + " dial err:" + err.Error())
	}
	fmt.Println("connect to server success:", cli.RemoteAddr().String(), resp.StatusCode)
	if resp == nil || resp.StatusCode != http.StatusSwitchingProtocols {
		return nil, fmt.Errorf("dial resp err:%v", resp.StatusCode)
	}
	ws := &webSocket{ctx: ctx, dialer: dialer, conn: cli, channel: make(chan []byte), err: make(chan error)}
	if method != "" {
		if err := cli.WriteMessage(subscribeType, []byte(method)); err != nil {
			return nil, errors.New("send method err:" + err.Error())
		}
		typ, resp, err := cli.ReadMessage()
		if err != nil && typ == subscribeType {
			ws.resp = string(resp)
		}
	}
	return ws, nil
}

func (ws *webSocket) Start(ctx context.Context) {
	go func() {
		fmt.Println("start receive routine")
		for {
			select {
			default:
				typ, body, err := ws.conn.ReadMessage()
				if err != nil {
					ws.err <- err
					return
				}
				fmt.Println("receive msg type:", typ)
				ws.channel <- body
			case <-ctx.Done():
				return
			}
		}
	}()

}

func (ws *webSocket) Err() chan error {
	return ws.err
}

func (ws *webSocket) Subscribe() chan []byte {
	return ws.channel
}

func (ws *webSocket) UnsubscribeMessage() string {
	return ws.resp
}
