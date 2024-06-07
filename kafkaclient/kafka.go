package kafkaclient

import (
	"bom-import-xls/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type Message struct {
	BomCode    string
	Components domain.BoM
}

func Publish[T any](message T, topic string) {

	//sarama.Logger = log.New(os.Stderr, "[KAFKA] ", log.LstdFlags)

	kafkaConn := sarama.NewConfig()
	kafkaConn.ClientID = "bom-import-xls"
	kafkaConn.Producer.Return.Successes = true
	kafkaConn.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConn.Producer.Retry.Max = 5
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, kafkaConn)
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

func Consume(ctx context.Context, topic string, consumer Consumer) (message []byte) {
	keepRunning := true
	brokers := "localhost:9092"
	//sarama.Logger = log.New(os.Stderr, "[KAFKA] ", log.LstdFlags)

	ctx, cancel := context.WithCancel(context.Background())
	kafkaConn := sarama.NewConfig()
	kafkaConn.Consumer.Return.Errors = true
	kafkaConn.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), "bom-import", kafkaConn)
	if err != nil {
		fmt.Println("Error consumer group: ", err)
	}

	//consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, kafkaConn)
	//if err != nil {
	//	fmt.Println("Error consumer: ", err)
	//}

	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	//defer func(consumer sarama.Consumer) {
	//	err := consumer.Close()
	//	if err != nil {
	//		fmt.Println("Error closing consumer: ", err)
	//	}
	//}(consumer)

	//partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	//if err != nil {
	//	fmt.Println("Error partition consumer: ", err)
	//}
	//
	//defer partitionConsumer.Close()
	//
	//for {
	//	select {
	//	case msg := <-partitionConsumer.Messages():
	//		fmt.Println("Received message: ", string(msg.Value))
	//		return msg.Value
	//
	//		// unmarshal message to BomMap
	//		//var bomMap domain.BomMap
	//		//err := json.Unmarshal(msg.Value, &bomMap)
	//		//if err != nil {
	//		//	fmt.Println("Error unmarshalling message: ", err)
	//		//}
	//		//
	//		//fmt.Println(bomMap)
	//		//ack
	//
	//	case err := <-partitionConsumer.Errors():
	//		fmt.Println("Error receiving message: ", err)
	//	}
	//}
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, strings.Split(topic, ","), &consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.Ready = make(chan bool)
		}
	}()
	<-consumer.Ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(client, &consumptionIsPaused)
		}
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
	return nil
}

type Consumer struct {
	Ready chan bool
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.Ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		log.Println("Resuming consumption")
	} else {
		client.PauseAll()
		log.Println("Pausing consumption")
	}

	*isPaused = !*isPaused
}
