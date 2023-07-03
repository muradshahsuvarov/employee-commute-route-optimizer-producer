package Producer

import (
	"fmt"
	"log"
	"main/src/EscapeCharacter"
	"main/src/Response"

	"github.com/Shopify/sarama"
	"gopkg.in/ini.v1"
)

func ProduceMessage(_id string, _server string, _topic string, _messageType string, _message string, _propertiesFile string) <-chan Response.Response {

	// Response Channel
	var responseChan chan Response.Response = make(chan Response.Response, 1)

	// Load producer properties from file
	filePath := _propertiesFile
	cfg, err := ini.Load(filePath)
	if err != nil {
		log.Printf("Failed to load producer properties file: %v", err)
		var res Response.Response = Response.Response{
			Id:        _id,
			Topic:     _topic,
			Partition: 0,
			Offset:    0,
			Error:     true,
		}
		responseChan <- res
		return responseChan
	}

	// Create Kafka config using the loaded properties
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = cfg.Section("producer").Key("retry.max").MustInt(5)
	config.Producer.Return.Successes = true // Set to true for SyncProducer

	fmt.Println("MAMA 0")
	// Create a new Kafka producer using the config
	producer, err := sarama.NewSyncProducer([]string{_server}, config)
	if err != nil {
		log.Printf("Failed to create Kafka producer: %v", err)
		var res Response.Response = Response.Response{
			Id:        _id,
			Topic:     _topic,
			Partition: 0,
			Offset:    0,
			Error:     true,
		}
		responseChan <- res
		return responseChan
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf("Failed to close Kafka producer: %v", err)
		}
	}()

	var escapedMessage string = EscapeCharacter.EscapeSpecialCharacters(_message)

	fmt.Println("MAMA 1")
	// Send a message to a Kafka topic
	topic := _topic
	message := fmt.Sprintf(`{"id":"%s","type":"%s","message":"%s"}`, _id, _messageType, escapedMessage)
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send Kafka message: %v", err)
		var res Response.Response = Response.Response{
			Id:        _id,
			Topic:     topic,
			Partition: partition,
			Offset:    offset,
			Error:     true,
		}
		responseChan <- res
		return responseChan
	}

	// Print the message details
	fmt.Printf("Message sent to topic '%s', partition %d, offset %d\n", topic, partition, offset)
	var res Response.Response = Response.Response{
		Id:        _id,
		Topic:     topic,
		Partition: partition,
		Offset:    offset,
		Error:     false,
	}
	responseChan <- res
	fmt.Println("MAMA 2")
	return responseChan
}
