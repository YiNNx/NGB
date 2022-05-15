package controller

import (
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"net/http"
	"ngb/model"
	"ngb/util"
	"strconv"
	"sync"
	"time"
)

//--------- Http Api ---------

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

// --------------- WebSocket Api --------------

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
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
		broadcast:  make(chan *model.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
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
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	from int
	to   int
	send chan *model.Message
}

func (c *Client) setID() string {
	return strconv.Itoa(c.from) + "_" + strconv.Itoa(c.to)
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		tx := model.BeginTx()
		_, content, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				util.Logger.Error(err)
			}
			util.Logger.Info(err)
			break
		}
		var rec receiveMessage
		err = json.Unmarshal(content, &rec)
		if err != nil {
			util.Logger.Error(err)
		}
		msg := &model.Message{
			Time:     time.Now(),
			Sender:   c.from,
			Receiver: c.to,
			Content:  rec.Content,
		}

		if err := model.Insert(msg); err != nil {
			tx.Rollback()
			util.Logger.Error(err)
		}

		hub.broadcast <- msg

		n := &model.Notification{
			Uid:       c.to,
			Type:      model.TypeMessage,
			ContentId: msg.Mid,
		}
		if err := model.Insert(n); err != nil {
			tx.Rollback()
			util.Logger.Error(err)
		}
		tx.Close()
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		c.conn.SetWriteDeadline(time.Now().Add(writeWait))
		msg, ok := <-c.send
		if !ok {
			// The hub closed the channel.
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		err := c.conn.WriteJSON(msg)
		if err != nil {
			util.Logger.Error(err)
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func getClient(from int, to int, ws *websocket.Conn) *Client {
	chatID := strconv.Itoa(from) + "_" + strconv.Itoa(to)
	if Client, ok := hub.clients[chatID]; ok {
		return Client
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

var wg sync.WaitGroup

func Chat(c echo.Context) error {

	with, err := strconv.Atoi(c.QueryParam("with"))
	if err != nil {
		util.Logger.Error(err)
		return err
	}
	uid := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		util.Logger.Error(err)
		return err
	}

	readClient := getClient(uid, with, ws)
	writeClient := getClient(with, uid, ws)

	go readClient.readPump()
	go writeClient.writePump()

	wg.Add(1)
	wg.Wait()
	return nil
}

// --------------- WebSocket Api --------------
//
//type Hub struct {
//	clients    map[string]*ClientChan
//	register   chan *ClientChan
//	unregister chan *ClientChan
//}
//
//func (hub *Hub) run() {
//	for {
//		select {
//		case client := <-hub.register:
//			chatID := setChatID(client.from, client.to)
//			hub.clients[chatID] = client
//			util.Logger.Debug(hub.clients)
//		case client := <-hub.unregister:
//			chatID := setChatID(client.from, client.to)
//			if _, ok := hub.clients[chatID]; ok {
//				delete(hub.clients, chatID)
//				close(client.msgChan)
//				util.Logger.Debug(hub.clients)
//			}
//		}
//	}
//}
//
//type ClientChan struct {
//	from    int
//	to      int
//	msgChan chan model.Message
//}
//
//func (client *ClientChan) close() {
//	util.Logger.Debug("close")
//
//	if _, ok := <-client.msgChan; ok {
//		close(client.msgChan)
//		util.Logger.Debug("close")
//	}
//	util.Logger.Debug("close")
//	hub.unregister <- client
//}
//
//func setChatID(from int, to int) string {
//	return strconv.Itoa(from) + "_" + strconv.Itoa(to)
//}
//
//func getClient(from int, to int, ws *websocket.Conn) *ClientChan {
//	chatID := setChatID(from, to)
//	if Client, ok := hub.clients[chatID]; ok {
//		return Client
//	}
//	client := &ClientChan{
//		from:    from,
//		to:      to,
//		msgChan: make(chan model.Message, 100),
//	}
//	hub.register <- client
//	return client
//}
//
//var upGrader = websocket.Upgrader{
//	ReadBufferSize:  1024,
//	WriteBufferSize: 1024,
//	CheckOrigin: func(r *http.Request) bool {
//		return true
//	},
//}
//
//func Read(ws *websocket.Conn, readChan *ClientChan) {
//	defer func() {
//		//hub.unregister <- readChan
//	}()
//
//	for {
//		tx := model.BeginTx()
//		_, content, err := ws.ReadMessage()
//		if err != nil {
//			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
//				util.Logger.Error(err)
//			}
//			util.Logger.Info(err)
//			return
//		}
//		var rec receiveMessage
//		err = json.Unmarshal(content, &rec)
//		if err != nil {
//			util.Logger.Error(err)
//		}
//		msg := model.Message{
//			Time:     time.Now(),
//			Sender:   readChan.from,
//			Receiver: readChan.to,
//			Content:  rec.Content,
//		}
//
//		if err := model.Insert(&msg); err != nil {
//			tx.Rollback()
//			util.Logger.Error(err)
//		}
//
//		readChan.msgChan <- msg
//
//		n := &model.Notification{
//			Uid:       readChan.to,
//			Type:      model.TypeMessage,
//			ContentId: msg.Mid,
//		}
//		if err := model.Insert(n); err != nil {
//			tx.Rollback()
//			util.Logger.Error(err)
//		}
//		tx.Close()
//	}
//}
//
//func Send(ws *websocket.Conn, c *ClientChan) {
//	defer func() {
//		//hub.unregister <- c
//	}()
//
//	for {
//		msg, ok := <-c.msgChan
//		if !ok {
//			break
//		}
//		err := ws.WriteJSON(msg)
//		if err != nil {
//			util.Logger.Error(err)
//		}
//	}
//}
//
//var wg sync.WaitGroup
//
//func Chat(c echo.Context) error {
//
//	with, err := strconv.Atoi(c.QueryParam("with"))
//	if err != nil {
//		util.Logger.Error(err)
//		return err
//	}
//	uid := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id
//
//	ws, err := upGrader.Upgrade(c.Response(), c.Request(), nil)
//	if err != nil {
//		util.Logger.Error(err)
//		return err
//	}
//	defer func(ws *websocket.Conn) {
//		err := ws.Close()
//		if err != nil {
//			util.Logger.Error("ws close error: " + err.Error())
//		}
//	}(ws)
//
//	readChan := getClient(uid, with, ws)
//	sendChan := getClient(with, uid, ws)
//
//	go Read(ws, readChan)
//	go Send(ws, sendChan)
//
//	wg.Add(1)
//	wg.Wait()
//	return nil
//}
