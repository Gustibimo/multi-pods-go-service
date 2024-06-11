package kafkaclient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
)

func Publish[T any](ctx context.Context, message T, topic string) {

	//sarama.Logger = log.New(os.Stderr, "[KAFKA] ", log.LstdFlags)

	kafkaConn := InitConfig()
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, kafkaConn.Configuration)
	if err != nil {
		fmt.Println("Error producer: ", err)
	}

	defer func(producer sarama.SyncProducer) {
		err := producer.Close()
		if err != nil {
			fmt.Println("Error closing producer: ", err)
		}
	}(producer)

	messageBytes, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshalling message: ", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(messageBytes),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("Error producer: ", err)
	}
	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
}
