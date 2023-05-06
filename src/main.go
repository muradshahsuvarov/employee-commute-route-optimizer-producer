package main

import (
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"gopkg.in/ini.v1"
)

func ProduceMessage(_server string, _topic string, _messageType string, _message string) (string, int32, int64) {
	// Load producer properties from file
	filePath := "<Path to producer.properties>"
	cfg, err := ini.Load(filePath)
	if err != nil {
		log.Fatalf("Failed to load producer properties file: %v", err)
	}

	// Create Kafka config using the loaded properties
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = cfg.Section("producer").Key("retry.max").MustInt(5)
	config.Producer.Return.Successes = true // Set to true for SyncProducer

	// Create a new Kafka producer using the config
	producer, err := sarama.NewSyncProducer([]string{_server}, config)
	if err != nil {
		log.Printf("Failed to create Kafka producer: %v", err)
		return _topic, 0, 0
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf("Failed to close Kafka producer: %v", err)
		}
	}()

	// Send a message to a Kafka topic
	topic := _topic
	message := "{\"type\":" + _messageType + ",\"message\":" + _message
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send Kafka message: %v", err)
		return topic, partition, offset
	}

	// Print the message details
	fmt.Printf("Message sent to topic '%s', partition %d, offset %d\n", topic, partition, offset)
	return topic, partition, offset
}

func main() {
	fmt.Println("Welcome to ECRO Kafka Producer")
}
