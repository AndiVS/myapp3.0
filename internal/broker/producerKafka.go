package broker

import (
	"github.com/AndiVS/myapp3.0/internal/model"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
)

// StartKafkaProducer start kafka producer
func StartKafkaProducer() *kafka.Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "172.28.1.6"})
	if err != nil {
		log.Printf("Kafka producers err %v", err)
	}

	return p
}

// ProduceEvent func for producing events of kafka
func (k *Kafka) ProduceEvent(destination, command string, data interface{}) {
	msgKafka := MessageForBrokers{
		Destination: destination,
		Command:     command,
		Cat:         data.(model.Cat),
	}

	msg, err := msgKafka.MarshalBinary()
	if err != nil {
		log.Printf("kafka marshaling err %v", err)
	}

	err = produceKafkaMsg(msg, k.Producer, k.Topic)
	if err != nil {
		log.Printf("err in produceKafkaMsg %v", err)
	}
}

func produceKafkaMsg(value []byte, p *kafka.Producer, topic string) error {
	deliveryChan := make(chan kafka.Event, 10000)
	err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value},
		deliveryChan,
	)
	if err != nil {
		log.Println("unable to enqueue message ", value)
	}

	go func() {
		e := <-deliveryChan
		m := e.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
		} else {
			log.Printf("Delivered message to topic %s [%d] at offset %v\n",
				*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
		}
		close(deliveryChan)
	}()
	return err
}
