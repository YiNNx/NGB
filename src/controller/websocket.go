package controller

import (
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"net/http"
	"ngb/model"
	"ngb/util"
	"ngb/util/log"
	"strconv"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Chat(c echo.Context) error {

	with, err := strconv.Atoi(c.QueryParam("with"))
	if err != nil {
		log.Logger.Error(err)
		return err
	}
	uid := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Logger.Error(err)
		return err
	}

	readClient := GetClient(uid, with, ws)
	writeClient := GetClient(with, uid, ws)

	go readClient.ReadMsg()
	go writeClient.WriteMsg()

	for {
	}
}

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

		n := &model.Notification{
			Uid:       c.to,
			Type:      model.TypeMessage,
			ContentId: msg.Mid,
		}
		if err := model.Insert(n); err != nil {
			tx.Rollback()
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

//--------- Old Http Api ---------

func SendMessage(c echo.Context) error {
	tx := model.BeginTx()
	defer tx.Close()

	receiver, err := strconv.Atoi(c.QueryParam("receiver"))
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	sender := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id
	rec := new(receiveMessage)
	if err := c.Bind(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := validate.Struct(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	m := &model.Message{
		Sender:   sender,
		Receiver: receiver,
		Content:  rec.Content,
	}
	if err := model.Insert(m); err != nil {
		tx.Rollback()
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	n := &model.Notification{
		Uid:       receiver,
		Type:      model.TypeMessage,
		ContentId: m.Mid,
	}
	if err := model.Insert(n); err != nil {
		tx.Rollback()
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return util.SuccessRespond(c, http.StatusOK, m)
}
