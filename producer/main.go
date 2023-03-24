package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	kafka "github.com/segmentio/kafka-go"
)

func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

type Message struct {
	MessageID string `json:"message_id`
	Data interface{} `json:"data"`
}

func main() {
	kafkaURL := os.Getenv("kafkaURL")
	topic := os.Getenv("topic")
	writer := newKafkaWriter(kafkaURL, topic)
	defer writer.Close()
	fmt.Println("start producing ... !!")

	data := `{"key": "17439118|consolidation.created", "topic": "consolidation.events", "offset": 7969691, "headers": {}, "payload": {"gigs": [{"id": 27653953, "size": "large", "state": "scheduled", "msa_id": "29820", "timing": {"pickup_end": null, "delivery_end": "2023-03-12T12:55:53Z", "pickup_start": null, "delivery_start": null, "pickup_timezone": "America/Los_Angeles"}, "category": "shipment", "distance": 3.6654675290000003, "profile_id": 321751, "pickup_time": null, "estimated_time": 13, "pickup_location": {"zip": "89119", "state": "NV", "latitude": 36.0838937, "longitude": -115.1537166}, "delivery_deadline": "2023-03-12T12:55:53Z", "delivery_location": {"latitude": 36.1050076, "longitude": -115.1731596}, "delivery_certifications": [{"id": 4, "name": "trusted_driver", "applies_to": {}, "created_at": "2020-11-06T22:39:05Z", "updated_at": "2023-01-18T21:11:29Z", "description": "Driver is established on the Roadie platform", "friendly_label": null}]}], "state": "published", "title": "Delta LASDL18245 - 3 bags - Support ID: pVjZKBnW", "payout": 10.32, "bearing": 311, "latitude": 36.0838937, "route_sk": "f17d7bb8-c0a2-11ed-935e-3e7994bd60c7", "longitude": -115.1537166, "created_at": "2023-03-12T06:55:58Z", "updated_at": "2023-03-12T06:55:58Z", "external_id": "ec2a5764-185c-4917-be60-b781c4f55d6a", "pickup_time": null, "publish_time": null, "consolidation_id": 17439118, "delivery_deadline": "2023-03-12T12:55:53Z", "payout_adjustment": 0.0, "estimated_distance": 3.67, "consolidation_stops": [{"index": 1, "gig_id": 27653953, "category": "pickup"}, {"index": 2, "gig_id": 27653953, "category": "delivery"}], "estimated_drive_time": 13}, "partition": 3, "timestamp": "2023-03-12T06:55:58.997+00:00"}`
	var message map[string]interface{}
	json.Unmarshal([]byte(data), &message)

	// fmt.Println(message["data"].(string))

	// data, err := json.Marshal(message["data"].(string))
	// if err != nil {
	// 	fmt.Println(err)
	// }

	fmt.Println(data)

	for i := 0; ; i++ {
		key := fmt.Sprint(uuid.New())
		msg := kafka.Message{
			Key:   []byte(key),
			Value: []byte(fmt.Sprint(data)),
		}
		err := writer.WriteMessages(context.Background(), msg)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("produced", key)
		}
		time.Sleep(1 * time.Second)
	}
}
