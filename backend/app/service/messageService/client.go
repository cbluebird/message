package messageService

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"message/app/comm"
	"message/app/model"
	"message/app/service/userService"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 5012

	readBufferSize  = 1024 // 读缓冲区大小
	writeBufferSize = 1024 // 写缓冲区大小
)

var (
	newline  = []byte{'\n'}
	space    = []byte{' '}
	upgrader = websocket.Upgrader{
		ReadBufferSize:  readBufferSize,
		WriteBufferSize: writeBufferSize,
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Client struct {
	ClientId int `json:"client_id"`
	Conn     *websocket.Conn
	send     chan *BroadcastChan
	hub      *Hub
}

type IncomingMessage struct {
	MessageType int    `json:"type"`
	Content     string `json:"content"`
	To          int    `json:"to"`
	SendType    int    `json:"send_type"`
}

// WsServer 处理websocket请求
func WsServer(c *gin.Context, msgType MsgType, clientId int) (*Client, error) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, err
	}
	client := &Client{
		ClientId: clientId,
		hub:      ClientHub,
		Conn:     conn,
		send:     make(chan *BroadcastChan, 256),
	}
	client.hub.ClientRegister <- client

	// 连接成功返回消息
	data := map[string]int{"client_id": client.ClientId}
	if err := WriteMessage(conn, Success, Success.Msg(), data, msgType); err != nil {
		return nil, err
	}

	// 监听客户端发送的消息
	go client.WriteMessageHandler(msgType)
	go client.ReadMessageHandler()
	go client.Init()

	return client, nil
}

func (c *Client) ReadMessageHandler() {
	if c.Conn != nil {
		defer func() {
			c.hub.ClientUnregister <- c
			c.Conn.Close()
		}()

		c.Conn.SetReadLimit(maxMessageSize)
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		c.Conn.SetPongHandler(func(appData string) error {
			c.Conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})
		for {
			_, message, err := c.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("Client closed connection:", err)
				}
				break
			}
			message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
			var incoming IncomingMessage
			if err := json.Unmarshal(message, &incoming); err != nil {
				log.Println("Invalid JSON message:", err)
				return
			}
			if err = CreateMessage(&model.Message{
				SendType: incoming.SendType,
				SendTo:   incoming.To,
				From:     c.ClientId,
				Ctime:    time.Now(),
				Content:  incoming.Content,
				MsgType:  comm.MessageType(incoming.MessageType),
			}); err != nil {
				log.Println("Failed to create message in db:", err)
				return
			}
			c.hub.ClientBroadcast <- &BroadcastChan{Content: incoming.Content, Form: c.ClientId, To: incoming.To, MessageType: incoming.MessageType, SendType: incoming.SendType, Ctime: time.Now().Format("2006-01-02 15:04:05")}
		}
	}
}

// WriteMessageHandler 将消息从集线器发送到 websocket 连接
func (c *Client) WriteMessageHandler(msgtype MsgType) {
	if c.Conn != nil {
		ticker := time.NewTicker(pingPeriod)
		defer func() {
			ticker.Stop()
			if c.Conn != nil {
				c.Conn.Close()
			}
		}()

		for {
			select {
			case message, ok := <-c.send:
				c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
				if !ok {
					c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				c.Conn.SetWriteDeadline(time.Time{})
				WriteMessage(c.Conn, SendMsgSuccess, SendMsgSuccess.Msg(), message, msgtype)
			case <-ticker.C:
				c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}
}

func (c *Client) Init() {
	user, err := userService.GetUserById(c.ClientId)
	if err != nil {
		log.Println("Failed to get user by ID:", err)
		return
	}
	messages, err := GetMessageList(user.Id, user.SignOutTime)
	if err != nil {
		log.Println("Failed to get message list:", err)
		return
	}
	if len(messages) > 0 {
		for _, message := range messages {
			if err := WriteMessage(c.Conn, Success, Success.Msg(), BroadcastChan{
				MessageType: int(message.MsgType),
				Content:     message.Content,
				To:          message.SendTo,
				SendType:    message.SendType,
				Form:        message.From,
				Ctime:       message.Ctime.Format("2006-01-02 15:04:05"),
			}, Json); err != nil {
				log.Println("Failed to send message:", err)
				return
			}
		}
	}
}
