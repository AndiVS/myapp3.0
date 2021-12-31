package broker

import (
	"context"

	"github.com/AndiVS/myapp3.0/internal/repository"
	"github.com/jackc/pgx/v4"
	"github.com/streadway/amqp"

	"log"
)

func closeCh(ch *amqp.Channel) {
	err := ch.Close()
	if err != nil {
		log.Printf("err ch clsoe %v ", err)
	}
}

// ConsumeEvents consume rabit event
func (r *RabbitMQ) ConsumeEvents(catsCont interface{}) {
	ch, err := r.Connection.Channel()
	if err != nil {
		log.Printf("err in %v", err)
	}
	defer closeCh(ch)

	q, err := ch.QueueDeclare(
		r.QName,
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Printf("err in %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Printf("err in %v", err)
	}

	forever := make(chan bool)

	b := new(pgx.Batch)
	msgCount := 0
	for d := range msgs {
		msgCount++
		ProcessMessage(d.Body, b)
		if msgCount%1000 == 0 {
			tx, errb := catsCont.(*repository.Postgres).Pool.Begin(context.Background())
			if errb != nil {
				log.Printf("err in tx begin %v ", errb)
			}

			batchResults := tx.SendBatch(context.Background(), b)
			var qerr error
			var rows pgx.Rows
			for qerr == nil {
				rows, qerr = batchResults.Query()
				rows.Close()
			}
			err = tx.Commit(context.Background())
			if err != nil {
				log.Printf("err in tx commit %v ", err)
			}
			b = new(pgx.Batch)
		}
		err = d.Ack(false)
		if err != nil {
			log.Printf("err Ack %v ", err)
		}
	}
	log.Printf(" [*] Awaiting RPC requests")
	<-forever
}

// ProcessMessage process mess
func ProcessMessage(m []byte, catsCont *pgx.Batch) {
	msg := new(MessageForBrokers)
	err := msg.UnmarshalBinary(m)
	if err != nil {
		log.Printf("err processmessasge kafka unmarhaling  %v", err)
	}

	switch msg.Destination {
	case userString:
	case catString:
		switch msg.Command {
		case insertCommand:
			catsCont.Queue(
				"INSERT INTO cats (_id, name, type) VALUES ($1, $2, $3) RETURNING _id",
				msg.Cat.ID, msg.Cat.Name, msg.Cat.Type)
			log.Printf("cat with id %v successfully inserted ", msg.Cat.ID)
		case deleteCommand:
			catsCont.Queue("DELETE FROM cats WHERE _id = $1", msg.Cat.ID)
			log.Printf("cat with id %v deleted successfully ", msg.Cat.ID)
		case updateCommand:
			catsCont.Queue("UPDATE cats SET name = $2, type = $3 WHERE _id = $1",
				msg.Cat.ID, msg.Cat.Name, msg.Cat.Type)
			log.Printf("cat with id %v updated successfully ", msg.Cat.ID)
		}
	}
}
