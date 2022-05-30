package util

import (
	"github.com/jordan-wright/email"
)

var pool *email.Pool
var hub = newHub()

func init() {
	//log.InitLog()
	go hub.run()
	initEs()
	initEmail()
}
