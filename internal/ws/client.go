package ws

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	waitTime = time.Duration(time.Second * 10)
)

type WsClient struct {
	messages  chan []byte
	conn      *websocket.Conn
	hub       *WsHub
	userID    uuid.UUID
	closeOnce sync.Once
}

func (wc *WsClient) writePump() {
	defer wc.Close()

	for msg := range wc.messages {
		wc.conn.SetWriteDeadline(time.Now().Add(waitTime))
		err := wc.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			wc.conn.WriteMessage(websocket.CloseMessage, nil)
			wc.hub.unregister <- wc
			return
		}

	}
}

func (wc *WsClient) Close() {
	wc.closeOnce.Do(func() {
		close(wc.messages)
		wc.conn.Close()
	})
}

func ServeWS(hub *WsHub, c *gin.Context) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// l error instead of print
		fmt.Println(err.Error())
		return
	}
	userID := c.MustGet("userID").(uuid.UUID)

	client := &WsClient{hub: hub, conn: conn, userID: userID, messages: make(chan []byte)}
	client.hub.register <- client

	go client.writePump()

}
