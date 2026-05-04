package ws

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	waitTime = time.Duration(time.Second * 10)
)

type WsClient struct {
	send   chan []byte
	conn   *websocket.Conn
	hub    *WsHub
	userID uuid.UUID
}

func (wc *WsClient) Read() {
	for {
		msg := <-wc.send
		wc.conn.SetWriteDeadline(time.Now().Add(waitTime))
		err := wc.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			// log instead
			fmt.Println(err.Error())

		}
	}
}

func ServeWS(hub *WsHub, c *gin.Context) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// log error instead of print
		fmt.Println(err.Error())
	}
	userID := c.MustGet("userID").(uuid.UUID)

	client := &WsClient{hub: hub, conn: conn, userID: userID, send: make(chan []byte)}
	client.hub.register <- client

	go client.Read()

}
