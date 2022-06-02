package util

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"ngb/config"
	"ngb/model"
	"ngb/util/log"
	"time"
)

var mqURL = "amqp://" + config.C.Rabbitmq.User + ":" + config.C.Rabbitmq.Password + "@" + config.C.Rabbitmq.Host + ":" + config.C.Rabbitmq.Port + "/"

type Notification struct {
	Time      time.Time
	Uid       int
	Email     string
	Type      int
	ContentId int
	Status    int
}

type sendMQ struct {
	content    []byte
	routingKey string
}

var mqChan = make(chan *sendMQ, 100)

var switchType = map[int]string{
	model.TypeMessage:   "message",
	model.TypeComment:   "comment",
	model.TypeMentioned: "mentioned",
	model.TypeNewPost:   "new_post",
}

func PublishToMQ(n *Notification) error {
	n.Time = time.Now()
	nBytes, err := json.Marshal(n)
	if err != nil {
		return err
	}
	mqChan <- &sendMQ{
		content:    nBytes,
		routingKey: switchType[n.Type],
	}
	return nil
}

func publish() {
	conn, err := amqp.Dial(mqURL)
	if err != nil {
		log.Logger.Error(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Logger.Error(err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		config.C.Rabbitmq.ExchangeName, // name
		"direct",                       // type
		true,                           // durable
		false,                          // auto-deleted
		false,                          // internal
		false,                          // no-wait
		nil,                            // arguments
	)
	if err != nil {
		log.Logger.Error(err)
	}

	for {
		send := <-mqChan
		err = ch.Publish(
			config.C.Rabbitmq.ExchangeName, // exchange
			send.routingKey,                // routing key
			false,                          // mandatory
			false,                          // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        send.content,
			})
		if err != nil {
			log.Logger.Error(err)
		}
	}

}
