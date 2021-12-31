package broker

import (
	"github.com/AndiVS/myapp3.0/internal/model"
	"github.com/google/uuid"
	"github.com/streadway/amqp"

	"log"
)

// StartCh start chen for rabbit
func (r *RabbitMQ) StartCh() {
	ch, err := r.Connection.Channel()
	if err != nil {
		log.Printf("err in %v", err)
	}
	if err != nil {
		log.Printf("err in %v", err)
	}
	r.Channel = ch
}

// ProduceEvent func for producing events
func (r *RabbitMQ) ProduceEvent(destination, command string, data interface{}) {
	msgRabbit := MessageForBrokers{
		Destination: destination,
		Command:     command,
		Cat:         data.(model.Cat),
	}

	ms, err := msgRabbit.MarshalBinary()
	if err != nil {
		log.Printf("kafka marshaling err %v", err)
	}
	corrID := uuid.New().String()

	produceRabbitMsg(ms, r.Channel, r.QName, corrID)
}

func produceRabbitMsg(value []byte, ch *amqp.Channel, qname, corrID string) {
	err := ch.Publish(
		"",
		qname,
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrID,
			Body:          value,
		})
	if err != nil {
		log.Printf("err in %v", err)
	}
}
