package Consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"main/src/KafkaResponse"
	"sync"

	"github.com/Shopify/sarama"
	"gopkg.in/ini.v1"
)

var msg *sarama.ConsumerMessage

func ConsumeMessages(_id string, _type string, _server string, _topic string, _partition int32, _propertiesFile string) <-chan KafkaResponse.KafkaResponse {

	fmt.Println("RUSTAM 1")
	// Declare KafkaResponse Channel
	kResChan := make(chan KafkaResponse.KafkaResponse, 1)

	// Load consumer properties from file
	filePath := _propertiesFile
	cfg, err_1 := ini.Load(filePath)

	if err_1 != nil {
		log.Fatalf("Failed to load consumer properties file: %v", err_1)
	}

	// Create Kafka config using the loaded properties
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Version = sarama.V2_6_0_0 // Set the desired Kafka version

	// Update Kafka config with properties from the file
	config.Net.SASL.Enable = cfg.Section("consumer").Key("sasl.enable").MustBool(false)
	config.Net.SASL.User = cfg.Section("consumer").Key("sasl.username").String()
	config.Net.SASL.Password = cfg.Section("consumer").Key("sasl.password").String()

	fmt.Println("RUSTAM 2")
	// Create a new Kafka consumer using the config
	consumer, err_2 := sarama.NewConsumer([]string{_server}, config)
	fmt.Println("RUSTAM 3")
	if err_2 != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err_2)
	}
	defer func() {
		if err_3 := consumer.Close(); err_3 != nil {
			log.Printf("Failed to close Kafka consumer: %v", err_3)
		}
	}()

	fmt.Println("RUSTAM 4")
	// Create a new Kafka consumer partition for the topic
	fmt.Println("OFFSET: ", sarama.OffsetOldest)
	partitionConsumer, err_4 := consumer.ConsumePartition(_topic, _partition, sarama.OffsetOldest)

	if err_4 != nil {
		log.Fatalf("Failed to create Kafka partition consumer: %v", err_4)
	}

	var kafkaResponse KafkaResponse.KafkaResponse = KafkaResponse.KafkaResponse{}

	wg := sync.WaitGroup{}

	wg.Add(1)
	fmt.Println("RUSTAM 5")
	// Start consuming messages from the Kafka topic
	go func() {

		defer close(kResChan)

		for {
			select {
			case msg = <-partitionConsumer.Messages():
				fmt.Println("RUSTAM 6")
				// Process the received message
				fmt.Printf("Received message: Topic=%s, Partition=%d, Offset=%d, Key=%s, Value=%s\n",
					msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				json.Unmarshal(msg.Value, &kafkaResponse)
				fmt.Println("DATA:", kafkaResponse)
				fmt.Printf("kafkaResponse.Id = %s\n_id = %s", kafkaResponse.Id, _id)
				if kafkaResponse.Id == _id {
					fmt.Println("\nkafkaResponse.Id == _id: ")
					kResChan <- kafkaResponse
					partitionConsumer.Close()
					wg.Done()
					return
				}
			case err := <-partitionConsumer.Errors():
				// Handle consumer errors
				log.Printf("Error while consuming message: %v\n", err)
				return
			}
		}
	}()

	fmt.Println("RUSTAM 7")
	wg.Wait()
	fmt.Println("RUSTAM 8")
	return kResChan

}
