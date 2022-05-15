package util

import (
	"fmt"
	"github.com/jordan-wright/email"
	"math/rand"
	"net/smtp"
	"net/textproto"
	"ngb/config"
	"strconv"
	"time"
)

var (
	Host     string = config.C.Mail.Host
	Addr     string = config.C.Mail.Addr
	Username string = config.C.Mail.Username
	Password string = config.C.Mail.Password

	coroutineNum int = 10
)

func SendEmail() {

	e := &email.Email{
		To:      []string{"2436201947@qq.com"},
		From:    Username,
		Subject: "Email Send Test",
		Text:    []byte("Text Body is, of course, supported!"),
		HTML:    []byte("<h1>This a test email</h1>"),
		Headers: textproto.MIMEHeader{},
	}

	err := e.Send(Addr, smtp.PlainAuth("", Username, Password, Host))
	if err != nil {
		Logger.Error(err)
	}
}

var ch chan *email.Email

var EmailList []string

func PushToPool() {
	if EmailList == nil {
		return
	}
	for i, _ := range EmailList {
		e := &email.Email{
			To:      []string{EmailList[i]},
			From:    Username,
			Subject: "Email Send Test",
			Text:    []byte("Text Body is, of course, supported!"),
			HTML:    []byte("<h1>This a test email</h1>"),
			Headers: textproto.MIMEHeader{},
		}
		ch <- e
	}
}

func Send(p *email.Pool) {
	for e := range ch {
		err := p.Send(e, 10*time.Second)
		if err != nil {
			Logger.Error(err)
		}
	}
}

func EmailPool() {
	p, err := email.NewPool(
		Addr,
		4,
		smtp.PlainAuth("", Username, Password, Host),
	)
	if err != nil {
		Logger.Error(err)
	}
	for i := 0; i < 4; i++ {
		go PushToPool()
		go Send(p)
	}
}

//-------------------------------------------

type Message struct {
	Id   int
	Name string
}

func Test2() {
	messages := make(chan Message, 100)
	result := make(chan error, 100)

	// 创建任务处理Worker
	for i := 0; i < coroutineNum; i++ {
		go worker(i, messages, result)
	}

	total := 0
	// 发布任务
	for k := 1; k <= 100; k++ {
		messages <- Message{Id: k, Name: "job" + strconv.Itoa(k)}
		total += 1
	}

	close(messages)

	// 接收任务处理结果
	for j := 1; j <= total; j++ {
		res := <-result
		if res != nil {
			fmt.Println(res.Error())
		}
	}

	close(result)
}

func worker(worker int, msg chan Message, result chan error) {
	// 从通道 chan Message 中监听&接收新的任务
	for job := range msg {
		fmt.Println("worker:", worker, "msg: ", job.Id, ":", job.Name)

		time.Sleep(time.Second * time.Duration(RandInt(1, 3)))

		result <- nil
	}
}

func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min+1)
}
