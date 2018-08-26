package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type room struct {
	// forwardは他のクライアントに転送するためのメッセージを保持するチャネルです
	forward chan [] byte
	// joinはチャットルームに参加しようとしているクライアントのためのチャネルです。
	join chan *client
	// leaveはチャットルームから退室しようとしているクライアントのためのチャネルです。
	leave chan *client
	// clientsには在室している全てのクライアントが保持されます
	clients map[*client]bool
}
func (r *room) run() {
	for {
		select {
		case client :=  <-r.join:
			// 参加
			r.clients[client]  = true
			case client := <-r.leave:
				// 退室
				delete(r.clients,client)
			close(client.send)
			case msg := <-r.forward:
				// 全てのクライアントにメッセージを転送
				for client := range r.clients {
					select {
					case client.send <- msg:
						// メッセージを送信
					default:
						//送信に失敗
						delete(r.clients,client)
					close(client.send)
					}
				}
		}
	}
}

const (
	socketBufferSize = 1024
	messageBufferSize = 256
	)
var upgrader = &websocket.Upgrader{ReadBufferSize:socketBufferSize,WriteBufferSize:socketBufferSize}

func (r* room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket,err := upgrader.Upgrade(w, req,nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client {
		socket:socket,
		send: make(chan []byte, messageBufferSize),
		room:r,
	}
	r.join <- client
	defer  func() {r.leave<-client}()
	go client.write()
	client.read()
}