package broker

import (
	"github.com/AndiVS/myapp3.0/internal/model"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
)

func StartKafkaConsumer() *kafka.Consumer {

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Panic("Kafka consumers err %v", err)
	}

	return c
}

func (k *Kafka) ConsumeEvents(catsMap map[string]*model.Cat) {
	err := k.Consumer.SubscribeTopics([]string{k.Topic}, nil)
	if err != nil {
		log.Printf("err in consumer subscribeTopic %v", err)
	}

	for {
		msg, err := k.Consumer.ReadMessage(-1)
		if err == nil {
			processMessage(msg, catsMap)
		} else {
			// The client will automatically try to recover from all errors.
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	/*	err = k.Consumer.Close()
		if err != nil {
			log.Printf("err in consumer close %v", err)
		}*/

}

func processMessage(msg *kafka.Message, catsMap map[string]*model.Cat) {
	msgKafka := new(MessageKafka)
	err := msgKafka.UnmarshalBinary(msg.Value)
	if err != nil {
		log.Printf("err processmessasge kafka unmarhaling  %v", err)
	}

	switch msgKafka.Destination {
	case "cat":
		switch msgKafka.Command {
		case "Insert":
			catsMap[msgKafka.Cat.ID.String()] = &msgKafka.Cat
			log.Printf("cat with id %v successfully inserted ", msgKafka.Cat.ID)
		case "Delete":
			delete(catsMap, msgKafka.Cat.ID.String())
			log.Printf("cat with id %v deleted successfully ", msgKafka.Cat.ID)
		case "Update":
			catsMap[msgKafka.Cat.ID.String()] = &msgKafka.Cat
			log.Printf("cat with id %v updated successfully ", msgKafka.Cat.ID)
		}
	}
}
