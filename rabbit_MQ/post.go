package rabbit_MQ

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func PostM1(body []byte, contType, id string) error {

	var err error
	for attempt := 0; attempt < attempts; attempt++ {
		_, err = c.ch.QueueDeclare(
			id,
			false,
			false,
			false,
			false,
			nil,
		)
		if err == nil {
			break
		}
	}
	failOnError(err, "failed to declare a queue")

	err = c.ch.PublishWithContext(c.ctx,
		"",
		c.qM1.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: contType,
			Body:        body,
		})

	if err != nil {
		return fmt.Errorf("failed to publish a message id[%s]: %s\n", id, err)
	}

	log.Printf("[%s]-> (Sent) %s\n", id, body[LenID:])
	return err
}

func PostM2(body []byte, contType, id string) error {

	err := c.ch.PublishWithContext(c.ctx,
		"",
		id,
		false,
		false,
		amqp.Publishing{
			ContentType: contType,
			Body:        body,
		})

	if err != nil {
		return fmt.Errorf("failed to publish a message id[%s]: %s\n", id, err)
	}
	return err
}
