package kafkaclient

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
)

func GetClient() sarama.Client {
	config := InitConfig()
	client, err := sarama.NewClient(config.BootstrapServers, config.Configuration)
	if err != nil {
		fmt.Println("Error client: ", err)
	}
	return client
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
