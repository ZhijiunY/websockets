package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// 設定一個 WebSocket 連線的升級器（upgrader），
// 它是用在將傳入的 HTTP 請求轉換成一個持久的 WebSocket 連線。
var (
	// websocket.Upgrader: 這是一個結構體，來自於 WebSocket 库，
	// 用於將 HTTP 連接升級到 WebSocket 連接。
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
	clients ClientList

	sync.RWMutex
}

// 創建工廠函數
func NewManager() *Manager {
	return &Manager{
		clients: make(ClientList),
	}
}

func (m *Manager) serverWS(w http.ResponseWriter, r *http.Request) {
	log.Println("New Connection")

	// upgrade regular http connection into websockte
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Create New Client
	client := NewClient(conn, m)
	m.addClient(client)

	// start client process
	go client.readMessage()
	go client.writeMessage()
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
	}
}
