package broker

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
)

func StartKafkaProducer() *kafka.Producer {

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "172.28.1.6:9092"})
	if err != nil {
		log.Panic("Kafka producers err %v", err)
	}

	return p
}

func (k *Kafka) ProduceEvent(destination, command string, data interface{}, topic string) {
	/*msgKafka := MessageKafka{
		Destination: destination,
		Command:     command,
		Cat:         data.(model.Cat),
	}

	msg, err := msgKafka.MarshalBinary()
	if err != nil {
		log.Printf("kafka marshaling err %v", err)
	}

	err = produceKafkaMsg(msg, k.Producer, topic)
	if err != nil {
		log.Printf("err in produceKafkaMsg %v", err)
	}*/
}

func produceKafkaMsg(value []byte, p *kafka.Producer, topic string) error {

	err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	}, nil)

	if err != nil {
		log.Println("unable to enqueue message ", value)
	}

	event := <-p.Events()
	message := event.(*kafka.Message)
	if message.TopicPartition.Error != nil {
		log.Println("Delivery failed due to error ", message.TopicPartition.Error)
	} else {
		log.Println("Delivered message to offset " + message.TopicPartition.Offset.String() + " in partition " + message.TopicPartition.String())
	}

	p.Flush(15 * 1000)

	return err
}
