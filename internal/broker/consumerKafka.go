package broker

import (
	"context"
	"github.com/AndiVS/myapp3.0/internal/repository"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
)

func StartKafkaConsumer() *kafka.Consumer {

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        "172.28.1.6:9092",
		"group.id":                 "myGroup",
		"auto.offset.reset":        "EARLIEST",
		"go.events.channel.enable": true,
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
			catsCont.(*repository.Postgres).InsertCat(context.Background(), &msgKafka.Cat)
			log.Printf("cat with id %v successfully inserted ", msgKafka.Cat.ID)
		case "Delete":
			catsCont.(*repository.Postgres).DeleteCat(context.Background(), msgKafka.Cat.ID)
			log.Printf("cat with id %v deleted successfully ", msgKafka.Cat.ID)
		case "Update":
			catsCont.(*repository.Postgres).UpdateCat(context.Background(), &msgKafka.Cat)
			log.Printf("cat with id %v updated successfully ", msgKafka.Cat.ID)
		}
	}
}
