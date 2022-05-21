package util

import (
	"github.com/jordan-wright/email"
)

var pool *email.Pool

func init() {
	//log.InitLog()
	initEs()
	initEmail()
}
