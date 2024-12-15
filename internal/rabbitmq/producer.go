package rabbitmq

import (
	"encoding/json"
	"strconv"

	"github.com/streadway/amqp"
)

type Producer interface {
	PublishTransactionID(id int) error
}

type producer struct {
	r *RabbitMQ
}

func NewProducer(r *RabbitMQ) Producer {
	return &producer{r: r}
}

func (p *producer) PublishTransactionID(id int) error {
	body, _ := json.Marshal(map[string]string{
		"transaction_id": strconv.Itoa(id),
	})
	return p.r.channel.Publish(
		"",
		p.r.queue.Name,
		false,
		false,
		amqpMessage(body),
	)
}

func amqpMessage(body []byte) amqp.Publishing {
	return amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}
}
