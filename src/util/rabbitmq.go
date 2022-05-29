package util

import (
	"github.com/streadway/amqp"
	"ngb/config"
	"ngb/util/log"
)

var mqURL = "amqp://" + config.C.Rabbitmq.User + ":" + config.C.Rabbitmq.Password + "@" + config.C.Rabbitmq.Host + ":" + config.C.Rabbitmq.Port + "/"

func Public(content []byte, routingKey string) error {
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
