package controller

import "github.com/go-playground/validator/v10"

var validate = validator.New()

//var hub = &Hub{
//	clients:    map[string]*ClientChan{},
//	register:   make(chan *ClientChan, 100),
//	unregister: make(chan *ClientChan, 100),
//}
var hub = newHub()

func init() {
	go hub.run()
}
