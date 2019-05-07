package common

import (
	"fmt"
	"github.com/hhy5861/go-common/logger"
	"github.com/streadway/amqp"
)

type (
	RabbitMQ struct {
		conn       *amqp.Connection
		ISerialize ISerialize
	}

	Payload struct {
		Exchange     string
		ExchangeType string
		RoutingKey   string
		Body         interface{}
	}

	RabbitMQConfig struct {
		Host     string `yaml:"host"`
		UserName string `yaml:"userName"`
		Password string `yaml:"password"`
		Port     int    `yaml:"port"`
	}
)

var (
	rabbitMQ *RabbitMQ
)

func NewRabbitClient() *RabbitMQ {
	return rabbitMQ
}

func NewRabbitMQ(config *RabbitMQConfig, iSerialize ...ISerialize) *RabbitMQ {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d",
		config.UserName,
		config.Password,
		config.Host,
		config.Port)

	conn, err := amqp.Dial(url)
	if err != nil {
		logger.Fatal(err)
	}

	defer conn.Close()

	var is ISerialize
	if len(iSerialize) > 0 {
		is = iSerialize[0]
	}

	rabbitMQ = &RabbitMQ{
		conn:       conn,
		ISerialize: is,
	}

	return rabbitMQ
}

func (r *RabbitMQ) getChannel() *amqp.Channel {
	channel, err := r.conn.Channel()
	if err != nil {
		logger.Fatal(err)
	}

	return channel
}

func (r *RabbitMQ) Declare(exchange, exchangeType string, reliable bool) *RabbitMQ {
	err := r.getChannel().ExchangeDeclare(exchange, exchangeType, true, false, false, false, nil)
	if err != nil {
		logger.Error(err)
	}

	if reliable {
		if err := r.getChannel().Confirm(false); err != nil {
			logger.Error(err)
			return nil
		}

		confirms := r.getChannel().NotifyPublish(make(chan amqp.Confirmation, 1))
		logger.Info(confirms)
	}

	return r
}

func (r *RabbitMQ) Publisher(payload Payload) error {
	data, err := r.ISerialize.Serialize(payload.Body)
	if err != nil {
		return err
	}

	body := amqp.Publishing{
		Headers:         amqp.Table{},
		ContentType:     "text/plain",
		ContentEncoding: "",
		Body:            data,
		DeliveryMode:    amqp.Transient,
		Priority:        0,
	}

	err = r.getChannel().Publish(payload.Exchange, payload.RoutingKey, false, false, body)
	if err != nil {
		return err
	}

	return nil
}
