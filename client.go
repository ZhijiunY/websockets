package main

import (
	"log"

	"github.com/gorilla/websocket"
)

// like handler
// 建立和管理 WebSocket 連接中的單個客戶端

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	egress     chan []byte
}

// 工廠模式
func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
	}
}

func (c *Client) readMessage() {

	defer func() {
		c.manager.removeClient(c)
	}()

	for {
		messageType, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {

				log.Printf("error reading message: %v", err)
			}
			break
	}

	for wsclient := range c.manager.clients {
		wsclient.egress <- payload
	}

	log.Println(messageType)
	log.Println(string(payload))
}

func (c *Client) writeMessage() {
	defer func() {
		c.manager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("conection closed: ", err)
				}
			}
			return
		}

		if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("failed to send message: %v", err)
		}
		log.Println("sent message")
	}
}
