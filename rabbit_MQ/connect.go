package rabbit_MQ

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"sync"
	"time"
)

const LenID = 8
const attempts = 3

type connect struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	qM1  amqp.Queue
	ctx  context.Context
}

var c connect

func ServerUP(adr string, timeout int, exit chan struct{}, wg *sync.WaitGroup) {

	serverUP(adr, timeout, exit, wg)
}

func serverUP(adr string, timeout int, exit chan struct{}, wg *sync.WaitGroup) {

	conn, err := amqp.Dial(adr)
	failOnError(err, "failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "failed to open a channel")
	defer ch.Close()

	qM1, err := ch.QueueDeclare(
		"M1",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(timeout)*time.Second)
	defer cancel()

	c = connect{
		conn: conn,
		ch:   ch,
		qM1:  qM1,
		ctx:  ctx,
	}
	wg.Done()

	<-exit
}

func failOnError(err error, msg string) {

	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
