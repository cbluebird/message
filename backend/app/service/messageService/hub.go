package messageService

import (
	"log"
	"message/app/model"
	"message/app/service/groupUserService"
	"message/app/service/userService"
	"sync"
)

type BroadcastChan struct {
	MessageType int    `json:"type"`
	Content     string `json:"content"`
	To          int    `json:"to"`
	SendType    int    `json:"send_type"`
	Form        int    `json:"form"`
	Ctime       string `json:"ctime"`
}

type Hub struct {
	Clients          map[*Client]bool // 全部客户端列表 {*Client1: bool, *Client2: bool...}
	ClientSet        map[int]*Client
	ClientRegister   chan *Client        // 客户端连接处理
	ClientUnregister chan *Client        // 客户端断开连接处理
	ClientLock       sync.RWMutex        // 客户端列表读写锁
	ClientBroadcast  chan *BroadcastChan // 来自客户端的入站消息 {Name:"clientId", Msg:"msg"}
}

// NewHub 实例化
func NewHub() *Hub {
	return &Hub{
		Clients:          make(map[*Client]bool),
		ClientSet:        make(map[int]*Client, 1000),
		ClientRegister:   make(chan *Client),
		ClientUnregister: make(chan *Client),
		ClientBroadcast:  make(chan *BroadcastChan, 1000),
	}
}

var ClientHub *Hub

// Run run chan listener
func (m *Hub) Run() {
	for {
		select {
		case client := <-m.ClientRegister:
			m.handleClientRegister(client)

		case client := <-m.ClientUnregister:
			m.handleClientUnregister(client)
			close(client.send)

		case clients := <-m.ClientBroadcast:
			m.ClientBroadcastHandle(clients)
		}
	}
}

// handleClientRegister 客户端连接处理
func (m *Hub) handleClientRegister(client *Client) {
	m.ClientLock.Lock()
	defer m.ClientLock.Unlock()
	m.ClientSet[client.ClientId] = client
	m.Clients[client] = true
}

// handleClientUnregister 客户端断开连接处理
func (m *Hub) handleClientUnregister(client *Client) {
	m.ClientLock.Lock()
	if err := userService.UpdateSignOutTime(client.ClientId); err != nil {
		log.Println(err)
		return
	}
	if _, ok := m.Clients[client]; ok {
		delete(m.Clients, client)
		delete(m.ClientSet, client.ClientId)
	}
	m.ClientLock.Unlock()
}

// ClientBroadcastHandle 单客户端通道处理
func (m *Hub) ClientBroadcastHandle(message *BroadcastChan) {
	if message.SendType == 1 {
		_client := m.ClientSet[message.To]
		if _client != nil {
			select {
			case _client.send <- message:
				break
			default:
				close(_client.send)
				m.handleClientUnregister(_client)
			}
		}
	} else if message.SendType == 2 {
		var users []*model.GroupUser
		var err error
		if users, err = groupUserService.GetUsersInGroup(message.To); err != nil {
			return
		}
		for _, user := range users {
			_client := m.ClientSet[user.UserId]
			if _client != nil {
				select {
				case _client.send <- message:
					break
				default:
					close(_client.send)
					m.handleClientUnregister(_client)
				}
			}
		}
	}
}
