package broker

import (
	"context"
	"github.com/AndiVS/myapp3.0/internal/repository"
	"github.com/jackc/pgx/v4"
	"log"
)

func (r *RabbitMQ) ConsumeEvents(catsCont interface{}) {

	ch, err := r.Connection.Channel()
	if err != nil {
		log.Printf("err in %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"rpc_queue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		log.Printf("err in %v", err)
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
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

	go func() {
		for d := range msgs {
			msgCount += 1
			processMessageR(d.Body, b)
			if msgCount%1000 == 0 {
				tx, err := catsCont.(*repository.Postgres).Pool.Begin(context.Background())
				if err != nil {
					log.Printf("err in tx begin %v ", err)
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

			/*	err = ch.Publish(
					"",        // exchange
					d.ReplyTo, // routing key
					false,     // mandatory
					false,     // immediate
					amqp.Publishing{
						ContentType:   "text/plain",
						CorrelationId: d.CorrelationId,
						Body:          []byte(strconv.Itoa(response)),
					})
				if err != nil {
					log.Printf("err in %v", err)
				}*/

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Awaiting RPC requests")
	<-forever
}

func processMessageR(m []byte, catsCont *pgx.Batch) {
	msg := new(MessageForBrokers)
	err := msg.UnmarshalBinary(m)
	if err != nil {
		log.Printf("err processmessasge kafka unmarhaling  %v", err)
	}

	switch msg.Destination {
	case "cat":
		switch msg.Command {
		case "Insert":
			catsCont.Queue(
				"INSERT INTO cats (_id, name, type) VALUES ($1, $2, $3) RETURNING _id",
				msg.Cat.ID, msg.Cat.Name, msg.Cat.Type)
			log.Printf("cat with id %v successfully inserted ", msg.Cat.ID)
		case "Delete":
			catsCont.Queue("DELETE FROM cats WHERE _id = $1", msg.Cat.ID)
			log.Printf("cat with id %v deleted successfully ", msg.Cat.ID)
		case "Update":
			catsCont.Queue("UPDATE cats SET name = $2, type = $3 WHERE _id = $1",
				msg.Cat.ID, msg.Cat.Name, msg.Cat.Type)
			log.Printf("cat with id %v updated successfully ", msg.Cat.ID)
		}
	}
}
