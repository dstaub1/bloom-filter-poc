package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"github.com/redis/go-redis/v9"
	kafka "github.com/segmentio/kafka-go"
)

// type KafkaMessage struct {
// 	Name string `json:"name"`
// 	MessageID string `json:"message_id"`
// 	Data string `json:"data"`
//  }

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

func handleRedis(rdb *redis.Client, topic string, key string, value string){
	switch topic {
	case "payment.events":
		err := rdb.Set(ctx, key, value, 0).Err()
		if err != nil{
			log.Println(err)
		}else {
			log.Printf("redis updated for topic %v : %v : %v",topic,key,value)
		}
	}
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(val)
}

var (
	ctx = context.Background()
)

type KafkaMessage struct {
	Name string `json:"name"`
	MessageID string `json:"message_id"`
	Data string `json:"data"`
 }

func main() {
	kafkaURL := os.Getenv("kafkaURL")
	topic := os.Getenv("topic")
	groupID := os.Getenv("groupID")

	reader := getKafkaReader(kafkaURL, topic, groupID)

	defer reader.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		Password: "",
		DB: 0,
	})
	fmt.Println("Connected to Redis...")

	fmt.Println("start consuming ... !!")
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		handleRedis(rdb, m.Topic, string(m.Key), string(m.Value))
	}
}
