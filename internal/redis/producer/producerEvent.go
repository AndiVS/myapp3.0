// Package producerredisredisredisredis for redis
package producer

import (
	"fmt"

	"github.com/go-redis/redis/v7"
)

// GenerateEvent for redis
func GenerateEvent(destination, command string, data interface{}, client *redis.Client, StreamName string) {
	newID, err := produceMsg(map[string]interface{}{
		"destination": destination,
		"command":     command,
		"data":        data,
	}, client, StreamName)

	checkError(err, command, destination, newID)
}

func produceMsg(values map[string]interface{}, client *redis.Client, StreamName string) (string, error) {
	str, err := client.XAdd(&redis.XAddArgs{
		Stream: StreamName,
		Values: values,
	}).Result()

	return str, err
}

func checkError(err error, request, record, requestID string) {
	if err != nil {
		fmt.Printf("get error:%v\n", err)
	} else {
		fmt.Printf("add to stream comand:%v for record:%v rqestID:%v\n", request, record, requestID)
	}
}
