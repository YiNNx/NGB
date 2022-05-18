package util

import (
	"github.com/jordan-wright/email"
	"net/smtp"
	"net/textproto"
	"ngb/config"
	"sync"
	"time"
)

var (
	host     string = config.C.Mail.Host
	addr     string = config.C.Mail.Addr
	username string = config.C.Mail.Username
	password string = config.C.Mail.Password

	coroutine int = config.C.Mail.Goroutine
	total     int
	count     int

	wg sync.WaitGroup
)

func EmailPool(emailList []string, subject string, text string) ([]*Result, error) {
	var emailChan = make(chan *email.Email, 100)
	var resChan = make(chan *Result, 100)

	total = len(emailList)
	count = 0

	p, err := email.NewPool(addr, 1, smtp.PlainAuth("", username, password, host))
	if err != nil {
		return nil, err
	}

	wg.Add(1)
	go pushToPool(emailList, subject, text, emailChan)

	for i := 0; i < coroutine; i++ {
		wg.Add(1)
		go send(p, emailChan, resChan)
	}

	var res []*Result

	wg.Add(1)
	go handleRes(resChan, &res)

	wg.Wait()
	return res, nil
}

func pushToPool(emailList []string, subject string, text string, emailChan chan *email.Email) {
	defer wg.Done()
	if emailList == nil {
		return
	}
	for i, _ := range emailList {
		e := &email.Email{
			To:      []string{emailList[i]},
			From:    username,
			Subject: subject,
			Text:    []byte(text),
			Headers: textproto.MIMEHeader{},
		}
		emailChan <- e
	}
	close(emailChan)
}

type Result struct {
	Email  []string
	Time   time.Time
	Status bool
	err    error
	Error  string
}

func send(p *email.Pool, emailChan chan *email.Email, resChan chan *Result) {

	defer wg.Done()
	for {
		e, ok := <-emailChan
		if !ok {
			break
		}
		err := p.Send(e, 10*time.Second)
		res := &Result{
			Email: e.To,
			Time:  time.Now(),
			err:   err,
		}

		resChan <- res
		count += 1
		if count == total {
			close(resChan)
		}
	}
}

func handleRes(resChan chan *Result, res *[]*Result) {
	defer wg.Done()

	for {
		r, ok := <-resChan
		if !ok {
			break
		}
		if r.err != nil {
			r.Status = false
			r.Error = r.err.Error()
			Logger.Error(r.Email, "send error: ", r.Error)
		} else {
			r.Status = true
			Logger.Info(r.Email, "successfully send!")
		}

		*res = append(*res, r)
	}
}
