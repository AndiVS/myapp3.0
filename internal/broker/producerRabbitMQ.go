package broker

import (
	"github.com/AndiVS/myapp3.0/internal/model"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"log"
)

func (r *RabbitMQ) ProduceEvent(destination, command string, data interface{}) {
	ch, err := r.Connection.Channel()
	if err != nil {
		log.Printf("err in %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		r.QName, // name
		false,   // durable
		false,   // delete when unused
		true,    // exclusive
		false,   // noWait
		nil,     // arguments
	)
	if err != nil {
		log.Printf("err in %v", err)
	}

	/*msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Printf("err in %v", err)
	}*/

	corrId := uuid.New().String()
	msgRabbit := MessageForBrokers{
		Destination: destination,
		Command:     command,
		Cat:         data.(model.Cat),
	}

	ms, err := msgRabbit.MarshalBinary()
	if err != nil {
		log.Printf("kafka marshaling err %v", err)
	}

	err = ch.Publish(
		"",          // exchange
		"rpc_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          ms,
		})
	if err != nil {
		log.Printf("err in %v", err)
	}

	/*for d := range msgs {
		if corrId == d.CorrelationId {
			res, err = strconv.Atoi(string(d.Body))
			if err != nil {
				log.Printf("err in %v", err)
			}
			break
		}
	}*/

}
