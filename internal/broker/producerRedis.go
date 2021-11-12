// Package producerredisredisredisredis for redis
package broker

import (
	"fmt"
	"log"

	"github.com/go-redis/redis/v7"
)

// ProduceEvent for redis
func (r *Redis) ProduceEvent(destination, command string, data interface{}, StreamName string) {
	newID, err := produceRedisMsg(map[string]interface{}{
		"destination": destination,
		"command":     command,
		"data":        data,
	}, r.Client, StreamName)

	checkError(err, command, destination, newID)
}

func produceRedisMsg(values map[string]interface{}, client *redis.Client, StreamName string) (string, error) {

	str, err := client.XAdd(&redis.XAddArgs{
		Stream: StreamName,
		Values: values,
	}).Result()
	if err != nil {
		log.Printf("err in add in stream %v", err)
	}
	return str, err
}

func checkError(err error, request, record, requestID string) {
	if err != nil {
		fmt.Printf("get error:%v\n", err)
	} else {
		fmt.Printf("add to stream comand:%v for record:%v rqestID:%v\n", request, record, requestID)
	}
}
