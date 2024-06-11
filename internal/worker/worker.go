package worker

import (
	"bom-import-xls/internal/worker/handlers"
	"bom-import-xls/kafkaclient"
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/hashicorp/go-uuid"
	_ "github.com/hashicorp/go-uuid"
	"log"
	"os"
	"os/signal"
	"sync"
)

var (
	stopChan = make(chan struct{})
	wg       sync.WaitGroup
)

func Run() {
	wg.Add(1)
	defer wg.Done()
	select {
	case <-stopChan:
		return
	default:
		// run actors
		StartConsumers("bom-import-xls", []string{"file-bom-import-parsing",
			"test-hello", "file-parsing-succeed", "file-parsing-failed"})
	}

}

func Stop() {
	// stop all worker
	close(stopChan)
}

func Wait() {
	wg.Wait()
	fmt.Println("Worker stopped.")
}

type Consumer struct {
	groupID string
	topics  []string
}

func (w *Consumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (w *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (w *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	traceID, _ := uuid.GenerateUUID()
	ctx := context.WithValue(context.Background(), "trace_id", traceID)
	for msg := range claim.Messages() {
		switch claim.Topic() {
		case "file-bom-import-parsing":
			handlers.HandleParseBomFile(ctx, string(msg.Value))
		case "test-hello":
			handlers.HandleHello(ctx)
		case "file-parsing-succeed":
			handlers.HandleFileParseSucceed(ctx, string(msg.Value))
		case "file-parsing-failed":
			handlers.HandleFileParseFailed(ctx)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

func StartConsumers(groupID string, topics []string) error {
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	client := kafkaclient.GetClient()
	consumerGroup, err := sarama.NewConsumerGroupFromClient(groupID, client)
	if err != nil {
		return err
	}

	worker := Consumer{groupID: groupID, topics: topics}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			if err := consumerGroup.Consume(ctx, topics, &worker); err != nil {
				log.Fatalf("Error from worker: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt)
	<-sigterm
	cancel()
	return consumerGroup.Close()
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
