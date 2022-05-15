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

// --------------- WebSocket --------------

var wg sync.WaitGroup

type Manager struct {
	Clients map[string]*ChatHub
}

var manager = &Manager{
	Clients: map[string]*ChatHub{},
}

type ChatHub struct {
	From    int
	To      int
	MsgChan chan model.Message
}

func setChatID(from int, to int) string {
	return strconv.Itoa(from) + "_" + strconv.Itoa(to)
}

func getChatHub(from int, to int) *ChatHub {
	chatID := setChatID(from, to)
	if chatHub, ok := manager.Clients[chatID]; ok {
		return chatHub
	}
	chatHub := &ChatHub{
		From:    from,
		To:      to,
		MsgChan: make(chan model.Message, 100),
	}
	manager.Clients[chatID] = chatHub
	return chatHub
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Chat(c echo.Context) error {

	with, err := strconv.Atoi(c.QueryParam("with"))
	if err != nil {
		util.Logger.Error(err)
		return err
	}
	uid := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id
	readHub := getChatHub(uid, with)
	sendHub := getChatHub(with, uid)

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		util.Logger.Error(err)
		return err
	}

	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			util.Logger.Error("ws close error: " + err.Error())
		}
	}(ws)

	wg.Add(1)
	go Read(ws, readHub)
	go Send(ws, sendHub)
	wg.Wait()
	return nil
}

func Read(ws *websocket.Conn, hub *ChatHub) {
	for {
		tx := model.BeginTx()
		_, content, err := ws.ReadMessage()
		if err != nil {
			util.Logger.Error(err)
		}
		var rec receiveMessage
		err = json.Unmarshal(content, &rec)
		if err != nil {
			util.Logger.Error(err)
		}
		msg := model.Message{
			Time:     time.Now(),
			Sender:   hub.From,
			Receiver: hub.To,
			Content:  rec.Content,
		}
		hub.MsgChan <- msg

		if err := model.Insert(msg); err != nil {
			tx.Rollback()
			util.Logger.Error(err)
		}

		n := &model.Notification{
			Uid:       hub.To,
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

func Send(ws *websocket.Conn, hub *ChatHub) {
	for {
		msg, ok := <-hub.MsgChan
		if !ok {
			break
		}
		err := ws.WriteJSON(msg)
		if err != nil {
			util.Logger.Error(err)
		}
	}
}
