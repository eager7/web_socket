package main

import (
	"context"
	"fmt"
	"github.com/eager7/web_socket"
)

const url = "ws://47.244.14.81:8586"

func main() {
	fmt.Println("web socket example...")
	ctx, cancel := context.WithCancel(context.Background())
	ws, err := web_socket.NewWebSocket(ctx, url, `{"method":"eth_subscribe","params":["newHeads"],"id":1,"jsonrpc":"2.0"}`)
	if err != nil {
		panic(err)
	}
	ws.Start(ctx)
	defer cancel()
	for {
		select {
		case err := <-ws.Err():
			fmt.Println("err receive:", err)
			return
		case msg := <-ws.Subscribe():
			fmt.Println("receive msg:", string(msg))
		}
	}
}
