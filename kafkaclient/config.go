package kafkaclient

import "github.com/IBM/sarama"

func ClientConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.ClientID = "bom-import-xls"
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Version = sarama.V2_1_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	return config
}

type Config struct {
	BootstrapServers []string
	GroupID          string
	Configuration    *sarama.Config
}

func InitConfig() *Config {
	Brokers := []string{"localhost:9092"}
	groupID := "bom-import-xls"
	return &Config{
		BootstrapServers: Brokers,
		GroupID:          groupID,
		Configuration:    ClientConfig(),
	}
}
