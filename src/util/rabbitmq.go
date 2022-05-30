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
	Type      int
	ContentId int //私信为mid,关注人发帖和@为pid,评论为cid
	Status    int //0未读 1已读
}

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
	err = publish(nBytes, switchType[n.Type])
	if err != nil {
		return err
	}

	return nil
}

func publish(content []byte, routingKey string) error {
	conn, err := amqp.Dial(mqURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
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
		return err
	}

	err = ch.Publish(
		config.C.Rabbitmq.ExchangeName, // exchange
		routingKey,                     // routing key
		false,                          // mandatory
		false,                          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        content,
		})
	log.Logger.Debug(content)
	if err != nil {
		return err
	}
	return nil
}
