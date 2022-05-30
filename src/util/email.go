package util

import (
	"github.com/jordan-wright/email"
	"net/smtp"
	"net/textproto"
	"ngb/config"
	"ngb/model"
	"ngb/util/log"
	"sync"
	"time"
)

var (
	host     = config.C.Mail.Host
	addr     = config.C.Mail.Addr
	username = config.C.Mail.Username
	password = config.C.Mail.Password

	coroutine = config.C.Mail.Goroutine

	wg sync.WaitGroup

	emailHub = make(chan *sending, 100)
)

type sending struct {
	emailList []string

	subject string
	text    string

	resChan chan *Result
}

func initEmail() {
	var err error
	pool, err = email.NewPool(addr, 1, smtp.PlainAuth("", username, password, host))
	if err != nil {
		log.Logger.Error(err)
	}

	for i := 0; i < coroutine; i++ {
		go send()
	}

	failList, subject, text := model.RedisReadFailList()
	if failList != nil {
		res, err := EmailPool(failList, subject, text)
		if err != nil {
			log.Logger.Error(err)
		}
		for i, _ := range res {
			log.Logger.Info(res[i].Email, res[i].Status, res[i].Error)
		}
	}
}

func EmailPool(emailList []string, subject string, text string) ([]*Result, error) {
	resultList := make(map[string]bool)
	for i, _ := range emailList {
		resultList[emailList[i]] = false
	}
	defer handlePanic(resultList, subject, text)

	if emailList == nil {
		return nil, nil
	}

	emailSending := &sending{
		emailList: emailList,
		subject:   subject,
		text:      text,
		resChan:   make(chan *Result, 100),
	}

	wg.Add(1)
	go pushToPool(emailSending)

	var res []*Result
	wg.Add(1)
	go handleRes(&res, &resultList, emailSending.resChan)

	wg.Wait()
	return res, nil
}

func pushToPool(emailSending *sending) {
	defer wg.Done()
	emailHub <- emailSending
}

type Result struct {
	Email  []string
	Time   time.Time
	Status bool
	err    error
	Error  string
}

func send() {
	for {
		s := <-emailHub
		for i, _ := range s.emailList {
			e := &email.Email{
				To:      []string{s.emailList[i]},
				From:    username,
				Subject: s.subject,
				Text:    []byte(s.text),
				Headers: textproto.MIMEHeader{},
			}
			err := pool.Send(e, 10*time.Second)
			res := &Result{
				Email: e.To,
				Time:  time.Now(),
				err:   err,
			}
			s.resChan <- res
		}
		close(s.resChan)
	}
}

func handleRes(res *[]*Result, resultList *map[string]bool, resChan chan *Result) {
	defer wg.Done()

	for {
		r, ok := <-resChan
		if !ok {
			break
		}
		if r.err != nil {
			r.Status = false
			r.Error = r.err.Error()
			log.Logger.Error(r.Email, "send error: ", r.Error)
		} else {
			r.Status = true
			log.Logger.Info(r.Email, "successfully send!")
		}
		for i, _ := range r.Email {
			(*resultList)[r.Email[i]] = true
		}
		*res = append(*res, r)
	}
}

func handlePanic(sendList map[string]bool, subject string, text string) {
	var failList []string
	for i, _ := range sendList {
		if (sendList)[i] == false {
			log.Logger.Error("email hasn't been sent:", i)
			failList = append(failList, i)
		}
	}
	if failList != nil {
		model.RedisSetFailList(failList, subject, text)
	}
}
