package model

import (
	"github.com/go-redis/redis/v7"
)

type Redis struct {
	Client     *redis.Client
	StreamName string
}

func NewRedisClient(client *redis.Client, streamName string) *Redis {
	return &Redis{Client: client, StreamName: streamName}
}
