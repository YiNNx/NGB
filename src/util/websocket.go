package util

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	//"ngb/controller"
	"ngb/model"
	"ngb/util/log"
	"strconv"
	"time"
)

type Hub struct {
	// Registered clients.
	clients map[string]*Client

	// Inbound messages from the clients.
	broadcast chan *model.Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan *model.Message, 100),
		register:   make(chan *Client, 100),
		unregister: make(chan *Client, 100),
		clients:    map[string]*Client{},
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			chatID := client.setID()
			h.clients[chatID] = client
		case client := <-h.unregister:
			chatID := client.setID()
			if _, ok := h.clients[chatID]; ok {
				delete(h.clients, chatID)
				close(client.send)
			}
		case message := <-h.broadcast:
			chatID := strconv.Itoa(message.Sender) + "_" + strconv.Itoa(message.Receiver)
			client := h.clients[chatID]
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, chatID)
			}

		}
	}
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	// The websocket connection.
	conn *websocket.Conn

	from int
	to   int
	send chan *model.Message
}

func (c *Client) setID() string {
	return strconv.Itoa(c.from) + "_" + strconv.Itoa(c.to)
}

type Message struct {
	Content string
}

func (c *Client) ReadMsg() {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()
	for {
		tx := model.BeginTx()
		_, content, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				log.Logger.Error(err)
			}
			log.Logger.Info(err)
			break
		}
		var rec Message
		err = json.Unmarshal(content, &rec)
		if err != nil {
			log.Logger.Error(err)
		}
		msg := &model.Message{
			Time:     time.Now(),
			Sender:   c.from,
			Receiver: c.to,
			Content:  rec.Content,
		}

		if err := model.Insert(msg); err != nil {
			tx.Rollback()
			log.Logger.Error(err)
		}

		hub.broadcast <- msg

		n := &Notification{
			Uid:       c.to,
			Type:      model.TypeMessage,
			ContentId: msg.Mid,
		}
		if err := PublishToMQ(n); err != nil {
			log.Logger.Error(err)
		}
		tx.Close()
	}
}

func (c *Client) WriteMsg() {
	defer func() {
		c.conn.Close()
	}()
	for {
		msg, ok := <-c.send
		if !ok {
			// The hub closed the channel.
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		err := c.conn.WriteJSON(msg)
		if err != nil {
			log.Logger.Error(err)
		}
	}
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func GetClient(from int, to int, ws *websocket.Conn) *Client {
	chatID := strconv.Itoa(from) + "_" + strconv.Itoa(to)
	if res, ok := hub.clients[chatID]; ok {
		client := &Client{
			conn: ws,
			from: from,
			to:   to,
			send: res.send,
		}
		return client
	}
	client := &Client{
		conn: ws,
		from: from,
		to:   to,
		send: make(chan *model.Message, 100),
	}
	hub.register <- client
	return client
}
