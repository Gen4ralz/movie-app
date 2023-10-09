package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gen4ralz/movie-app/rating-service/pkg/model"
)

// Ingester defines a Kafka ingester.
type Ingester struct {
	consumer	*kafka.Consumer
	topic		string
}

// NewIngester creates a new Kafka ingester.
func NewIngester(addr string, groupID string, topic string) (*Ingester, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": addr,
		"group.id": groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}
	return &Ingester{
		consumer: consumer,
		topic:  topic,
		}, nil
}

// Ingest starts ingestion from Kafka and returns a channel containing rating events
// representing the data consumed from the topic.
func (i *Ingester) Ingest(ctx context.Context) (chan model.RatingEvent, error) {
	err := i.consumer.SubscribeTopics([]string{i.topic}, nil)
	if err != nil {
		return nil, err
	}

	ch := make(chan model.RatingEvent, 1)
	go func(){
		for {
			select {
			case <- ctx.Done():
				close(ch)
				i.consumer.Close()
				return
			default:
			}

			// -1 means always consume from the beginning of the topic,
			// reading all existing messages.
			msg, err := i.consumer.ReadMessage(-1)
			if err != nil {
				fmt.Println("Consumer error: " + err.Error())
				continue
			}

			var event model.RatingEvent
			err = json.Unmarshal(msg.Value, &event)
			if err != nil {
				fmt.Println("Unmarshal error: " + err.Error())
				continue
			}
			ch <- event
		}
	}()

	return ch, nil
}