package broker

import (
	"github.com/AndiVS/myapp3.0/internal/model"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-redis/redis/v7"
	"github.com/vmihailenco/msgpack/v4"
)

type Broker interface {
	ProduceEvent(destination, command string, cat *model.Cat, topic string)
	ConsumeEvents(catsMap map[string]*model.Cat)
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

type MessageKafka struct {
	Destination string    `json:"destination"`
	Command     string    `json:"command"`
	Cat         model.Cat `json:"cat"`
}

// MarshalBinary Marshal cat for redis stream
func (msg *MessageKafka) MarshalBinary() ([]byte, error) {
	return msgpack.Marshal(msg)
}

// UnmarshalBinary Marshal cat for redis stream
func (msg *MessageKafka) UnmarshalBinary(data []byte) error {
	return msgpack.Unmarshal(data, msg)
}
