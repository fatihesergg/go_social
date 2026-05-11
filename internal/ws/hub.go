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
	Messages   chan WsData
}

func NewWsHub() *WsHub {
	return &WsHub{
		Clients:    map[uuid.UUID]*WsClient{},
		register:   make(chan *WsClient),
		unregister: make(chan *WsClient),
		Messages:   make(chan WsData),
	}
}

func (ws *WsHub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			ws.Close()
			return
		case client := <-ws.register:
			ws.RegisterClient(client)
		case client := <-ws.unregister:
			ws.RemoveClient(client.userID)
		case message := <-ws.Messages:
			ws.SendClient(message)
		}
	}
}

func (ws *WsHub) RegisterClient(client *WsClient) {
	ws.Clients[client.userID] = client
}

func (ws *WsHub) GetClient(userID uuid.UUID) *WsClient {
	client, ok := ws.Clients[userID]
	if !ok {
		return nil
	}
	return client
}

func (ws *WsHub) RemoveClient(userID uuid.UUID) {
	_, ok := ws.Clients[userID]
	if !ok {
		return
	}
	delete(ws.Clients, userID)
}

func (ws *WsHub) SendClient(message WsData) {
	if client, ok := ws.Clients[message.UserID]; ok {
		client.messages <- message.Data
	}
}

func (ws *WsHub) Close() {
	for _, client := range ws.Clients {
		client.Close()

	}

}
