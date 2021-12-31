package broker

import (
	"context"

	"github.com/AndiVS/myapp3.0/internal/repository"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

// StartKafkaConsumer start kafka consumer
func StartKafkaConsumer() *kafka.Consumer {
	// group := os.Getenv("GROP")
	group := "gr1"
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "172.28.1.6",
		"group.id":          group,
		"auto.offset.reset": "latest",
	})
	if err != nil {
		log.Printf("Kafka consumers err %v", err)
	}

	return c
}

// ConsumeEvents consume kafka event
func (k *Kafka) ConsumeEvents(catsCont interface{}) {
	// err := k.Consumer.SubscribeTopics([]string{k.Topic}, nil)
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

	flag := true
	msgCount := 0
	for flag {
		ev := k.Consumer.Poll(0)
		switch e := ev.(type) {
		case *kafka.Message:
			msgCount++
			ProcessMessage(e.Value, b)
			if msgCount%1000 == 0 {
				tx, er := catsCont.(*repository.Postgres).Pool.Begin(context.Background())
				if er != nil {
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
					log.Printf("err in tx commit in postgres %v ", err)
				}
				b = new(pgx.Batch)
				go func() {
					_, errK := k.Consumer.Commit()
					if errK != nil {
						log.Printf("err in tx commit in kafka %v ", errK)
					}
				}()
			}
		case kafka.PartitionEOF:
			log.Printf("%% Reached %v\n", e)
		case kafka.Error:
			log.Printf("%% Error: %v\n", e)
			flag = false
		default:
			// fmt.Printf("Ignored %v\n", e)
		}
	}

	err = k.Consumer.Close()
	if err != nil {
		log.Printf("consumer close err  %v", err)
	}
}
