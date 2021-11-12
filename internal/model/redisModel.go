package model

import (
	"github.com/go-redis/redis/v7"
)

// Redis struct for redis
type Redis struct {
	Client     *redis.Client
	StreamName string
}

// NewRedisClient client for redis
func NewRedisClient(client *redis.Client, streamName string) *Redis {
	return &Redis{Client: client, StreamName: streamName}
}
