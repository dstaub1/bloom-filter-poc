package main

import (
   "context"
	"fmt"
   "encoding/json"
	"github.com/redis/go-redis/v9"
)

type KafkaMessage struct {
   Name string `json:"name"`
   MessageID string `json:"message_id"`
   Data string `json:"data"`
}

var (
         ctx = context.Background()
         rdb *redis.Client
)

func main() {
   rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

   // pong, err := rdb.Ping(ctx).Result()
	// fmt.Println(pong, err)

   json, err := json.Marshal(KafkaMessage{Name: "PaymentEventsConsumer", MessageID: "16cf77b4-d7af-4b0f-955b-420187f11943", Data: "Some json data"})
   if err != nil {
      fmt.Println(err)
   }

   err = rdb.Set(ctx, "abcdefg", json, 0).Err()
   if err != nil {
      fmt.Println(err)
   }
   val, err := rdb.Get(ctx, "abcdefg").Result()
   if err != nil {
      fmt.Println(err)
   }
   fmt.Println(val)

}

// var Pool *redis.Pool

// type Video struct{
//    Title string `redis:"title"`
//    Category string `redis:"category"`
//    Likes int `redis:"likes"`
// }

// func Init(){
//    Pool = &redis.Pool{
//       MaxIdle:     10,
//       IdleTimeout: 240 * time.Second,
//       Dial: func() (redis.Conn, error) {
//          return redis.Dial("tcp", "localhost:6379")
//       },
//    }
// }