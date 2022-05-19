package controller

import "github.com/go-playground/validator/v10"

var validate = validator.New()
var hub = newHub()

func init() {
	go hub.run()
}

