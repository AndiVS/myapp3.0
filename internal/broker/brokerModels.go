package broker

import (
	"encoding/json"
	"github.com/AndiVS/myapp3.0/internal/model"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-redis/redis/v7"
	"github.com/streadway/amqp"
	"log"
)

type Broker interface {
	ProduceEvent(destination, command string, data interface{})
	ConsumeEvents(interface{})
	GetString() string
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

func (r *Redis) GetString() string {
	return r.StreamName
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

func (k *Kafka) GetString() string {
	return k.Topic
}

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

type RabbitMQ struct {
	Connection *amqp.Connection
	QName      string
}

func (r *RabbitMQ) GetString() string {
	return r.QName
}

// NewRabbitMQ client for kafka
func NewRabbitMQ(qname string) Broker {
	conn, err := amqp.Dial("amqp://andeisaldyun:e3cr3t@172.28.1.7:5672/")
	if err != nil {
		log.Printf(" rabbit con err %v", err)
	}
	return &RabbitMQ{Connection: conn, QName: qname}
}
