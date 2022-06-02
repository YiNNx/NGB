package util

import (
	"github.com/jordan-wright/email"
)

var pool *email.Pool
var hub = newHub()

func init() {
	initEs()
	initEmail()

	go hub.run()
	go publish()
}
