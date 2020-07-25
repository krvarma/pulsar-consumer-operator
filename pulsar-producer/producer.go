package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

func getEnv(key string, defvalue string) string {
	value := os.Getenv(key)

	if value == "" {
		value = defvalue
	}

	return value
}

func main() {
	server := getEnv("PULSAR_SERVER", "pulsar://ubuntuserver:6650")
	message := getEnv("PULSAR_MESSAGE", "Sample Message213")
	topic := getEnv("PULSAR_TOPIC", "my-topic")

	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               server,
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})

	if err != nil {
		log.Fatalf("Could not instantiate Pulsar client: %v", err)
	}

	defer client.Close()

	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: topic,
	})

	_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: []byte(message),
	})

	defer producer.Close()

	if err != nil {
		fmt.Println("Failed to publish message", err)
	}

	fmt.Println("Published message")
}
