package rabbit_MQ

import (
	"github.com/rabbitmq/amqp091-go"
)

type Work struct {
	ContType string
	ID       string
	Body     []byte
}

func GetFromM1(w chan Work) {

	msgs, err := c.ch.Consume(
		c.qM1.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "failed to register a consumer")

	go func() {
		for d := range msgs {
			getData(d, w)
		}
	}()

	var forever chan struct{}
	<-forever
}

func GetFromM2(id string) (body []byte, cType string) {

	var err error
	var msgs <-chan amqp091.Delivery

	for attempt := 0; attempt < attempts; attempt++ {
		msgs, err = c.ch.Consume(
			id,
			"",
			true,
			false,
			false,
			false,
			nil,
		)
		if err == nil {
			break
		}
	}
	failOnError(err, "failed to register a consumer")

	d := <-msgs
	return d.Body, d.ContentType
}

func getData(d amqp091.Delivery, w chan Work) {

	w <- Work{
		ContType: d.ContentType,
		ID:       string(d.Body[:LenID]),
		Body:     d.Body[LenID:],
	}
}
