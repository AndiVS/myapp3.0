package broker

import (
	"github.com/AndiVS/myapp3.0/internal/model"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
)

func StartKafkaProducer() *kafka.Producer {

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "172.28.1.5"})
	if err != nil {
		log.Panic("Kafka producers err %v", err)
	}

	return p
}

func (k *Kafka) ProduceEvent(destination, command string, data interface{}, topic string) {
	msgKafka := MessageKafka{
		Destination: destination,
		Command:     command,
		Cat:         data.(model.Cat),
	}

	msg, err := msgKafka.MarshalBinary()
	if err != nil {
		log.Printf("kafka marshaling err %v", err)
	}

	go func() {
		for e := range k.Producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					log.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	err = produceKafkaMsg(msg, k.Producer, topic)
	if err != nil {
		log.Printf("err in produceKafkaMsg %v", err)
	}
}

func produceKafkaMsg(value []byte, p *kafka.Producer, topic string) error {

	err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	}, nil)

	p.Flush(15 * 1000)

	return err
}
