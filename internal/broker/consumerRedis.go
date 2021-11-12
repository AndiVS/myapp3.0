package broker

import (
	"github.com/AndiVS/myapp3.0/internal/model"
	"github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
)

// ConsumeEvents consume events
func (r *Redis) ConsumeEvents(catsMap map[string]*model.Cat) {
	for {
		streams, err := r.Client.XRead(&redis.XReadArgs{
			Streams: []string{r.StreamName, "$"},
		}).Result()
		if err != nil {
			log.Printf("err on consume events: %+v\n", err)
		}

		stream := streams[0].Messages[0]
		processStream(stream, catsMap)
	}
}

func processStream(stream redis.XMessage, catsMap map[string]*model.Cat) {
	destination := stream.Values["destination"].(string)
	command := stream.Values["command"].(string)

	switch destination {
	/*	case "user":
		user := new(model.User)
		err := msgpack.Unmarshal([]byte(stream.Values["data"].(string)),user)
		if err!= nil {}

		switch command {
		case "Update":
			err = r.ServiceUsers.UpdateUser(context.Background(), user.Username, user.IsAdmin)
			if err != nil {
				log.Printf("err %v", err )
			}else {
				log.Printf("user with username %v  successfully ", user.Username)
			}
		case "Delete":
			err = r.ServiceUsers.DeleteUser(context.Background(), user.Username)
			if err != nil {
				log.Printf("err %v", err )
			}else {
				log.Printf("user with username %v deleted successfully ", user.Username)
			}
		}*/
	case "cat":
		cat := new(model.Cat)

		err := cat.UnmarshalBinary([]byte(stream.Values["data"].(string)))
		if err != nil {
			log.Printf("err %v ", err)
		}

		switch command {
		case "Insert":
			catsMap[cat.ID.String()] = cat
			log.Printf("cat with id %v successfully inserted ", cat.ID)
		case "Delete":
			delete(catsMap, cat.ID.String())
			log.Printf("cat with id %v deleted successfully ", cat.ID)
		case "Update":
			catsMap[cat.ID.String()] = cat
			log.Printf("cat with id %v updated successfully ", cat.ID)
		}
	}
}
