package broker

import (
	"encoding/json"

	"github.com/AndiVS/myapp3.0/internal/model"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-redis/redis/v7"
	"github.com/streadway/amqp"

	"log"
)

const userString = "user"
const catString = "cat"
const updateCommand = "Update"
const deleteCommand = "Delete"
const insertCommand = "Insert"

// Broker interface for brokers
type Broker interface {
	ProduceEvent(destination, command string, data interface{})
	ConsumeEvents(interface{})
}

// Redis struct for redis
type Redis struct {
	Client     *redis.Client
	StreamName string
}

// NewRedisClient client for redis
func NewRedisClient(client *redis.Client, streamName string) Broker {
	return &Redis{Client: client, StreamName: streamName}
}

// Kafka struct for kafka
type Kafka struct {
	Consumer *kafka.Consumer
	Producer *kafka.Producer
	Topic    string
}

// NewKafka client for kafka
func NewKafka(consumer *kafka.Consumer, producer *kafka.Producer, topic string) Broker {
	return &Kafka{Consumer: consumer, Producer: producer, Topic: topic}
}

// MessageForBrokers message that send kafka and rabbit
type MessageForBrokers struct {
	Destination string    `param:"destination" query:"destination" header:"destination" form:"destination" json:"destination" xml:"destination" bson:"destination"`
	Command     string    `param:"command" query:"command" header:"command" form:"command" json:"command" xml:"command" bson:"command"`
	Cat         model.Cat `param:"cat" query:"cat" header:"cat" form:"cat" json:"cat" xml:"cat" bson:"cat"`
}

// MarshalBinary Marshal cat for redis stream
func (msg *MessageForBrokers) MarshalBinary() ([]byte, error) {
	return json.Marshal(msg)
}

// UnmarshalBinary Marshal cat for redis stream
func (msg *MessageForBrokers) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, msg)
}

// RabbitMQ struct for rabbitmq
type RabbitMQ struct {
	Connection *amqp.Connection
	QName      string
	Channel    *amqp.Channel
}

// NewRabbitMQ client for kafka
func NewRabbitMQ(qname string) Broker {
	// conn, err := amqp.Dial("amqp://andeisaldyun:e3cr3t@172.28.1.7:5672/")
	conn, err := amqp.Dial("amqp://andeisaldyun:e3cr3t@rabbitmq:5672/")
	if err != nil {
		log.Printf(" rabbit con err %v", err)
	}
	return &RabbitMQ{Connection: conn, QName: qname}
}
