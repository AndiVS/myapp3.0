package broker

import (
	"context"
	"github.com/AndiVS/myapp3.0/internal/repository"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
	"os"
)

func StartKafkaConsumer() *kafka.Consumer {

	group := os.Getenv("GROP")
	//group := "gr1"
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "172.28.1.6",
		"group.id":          group,
		"auto.offset.reset": "latest",
	})
	if err != nil {
		log.Panic("Kafka consumers err %v", err)
	}

	return c
}

func (k *Kafka) ConsumeEvents(catsCont interface{}) {
	//err := k.Consumer.SubscribeTopics([]string{k.Topic}, nil)
	err := k.Consumer.Subscribe(k.Topic, nil)
	if err != nil {
		log.Println("Unable to subscribe to topic " + k.Topic + " due to error - " + err.Error())
	} else {
		log.Println("subscribed to topic ", k.Topic)
	}
	/*
		for {
			log.Println("waiting for event...")
			kafkaEvent := <-k.Consumer.Events()
			if kafkaEvent != nil {
				switch event := kafkaEvent.(type) {
				case *kafka.Message:
					processMessage(event, catsCont)
				case kafka.Error:
					log.Println("Consumer error ", event.String())
				case kafka.PartitionEOF:
					log.Println(kafkaEvent)
				default:
					log.Println(kafkaEvent)
				}
			} else {
				log.Println("Event was null")
			}
		}*/

	b := new(pgx.Batch)

	run := true
	msgCount := 0
	for run == true {
		ev := k.Consumer.Poll(0)
		switch e := ev.(type) {
		case *kafka.Message:
			msgCount += 1
			processMessage(e, b)
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
				go func() {
					_, err := k.Consumer.Commit()
					if err != nil {
					}

				}()
			}
		case kafka.PartitionEOF:
			log.Printf("%% Reached %v\n", e)
		case kafka.Error:
			log.Printf("%% Error: %v\n", e)
			run = false
		default:
			//fmt.Printf("Ignored %v\n", e)
		}
	}

	err = k.Consumer.Close()
	if err != nil {
		log.Printf("consumer close err  %v", err)
	}
}

func processMessage(msg *kafka.Message, catsCont interface{}) {
	msgKafka := new(MessageKafka)
	err := msgKafka.UnmarshalBinary(msg.Value)
	if err != nil {
		log.Printf("err processmessasge kafka unmarhaling  %v", err)
	}

	switch msgKafka.Destination {
	case "cat":
		switch msgKafka.Command {
		case "Insert":
			catsCont.(*pgx.Batch).Queue(
				"INSERT INTO cats (_id, name, type) VALUES ($1, $2, $3) RETURNING _id",
				msgKafka.Cat.ID, msgKafka.Cat.Name, msgKafka.Cat.Type)
			log.Printf("cat with id %v successfully inserted ", msgKafka.Cat.ID)
		case "Delete":
			catsCont.(*pgx.Batch).Queue("DELETE FROM cats WHERE _id = $1", msgKafka.Cat.ID)
			log.Printf("cat with id %v deleted successfully ", msgKafka.Cat.ID)
		case "Update":
			catsCont.(*pgx.Batch).Queue("UPDATE cats SET name = $2, type = $3 WHERE _id = $1",
				msgKafka.Cat.ID, msgKafka.Cat.Name, msgKafka.Cat.Type)
			log.Printf("cat with id %v updated successfully ", msgKafka.Cat.ID)
		}
	}
}
