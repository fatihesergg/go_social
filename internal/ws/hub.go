package ws

import (
	"context"

	"github.com/google/uuid"
)

type WsData struct {
	Data   []byte
	UserID uuid.UUID
}

type WsHub struct {
	//NOTE: one connection per user for now
	Clients    map[uuid.UUID]*WsClient
	register   chan *WsClient
	unregister chan *WsClient
	Broadcast  chan WsData
}

func NewWsHub() *WsHub {
	return &WsHub{
		Clients:    map[uuid.UUID]*WsClient{},
		register:   make(chan *WsClient),
		unregister: make(chan *WsClient),
		Broadcast:  make(chan WsData),
	}
}

func (ws *WsHub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			ws.Close()
			return
		case client := <-ws.register:
			ws.Clients[client.userID] = client
		case client := <-ws.unregister:
			if _, ok := ws.Clients[client.userID]; ok {
				delete(ws.Clients, client.userID)
				close(client.send)
			}
		case message := <-ws.Broadcast:
			if client, ok := ws.Clients[message.UserID]; ok {
				client.send <- message.Data
			}
		}
	}
}

func (ws *WsHub) Close() {
	for _, client := range ws.Clients {
		close(client.send)
	}
	clear(ws.Clients)
	close(ws.Broadcast)
	close(ws.register)
	close(ws.unregister)
}
