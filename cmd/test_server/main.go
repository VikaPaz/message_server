package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/VikaPaz/message_server/internal/models"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

type Config struct {
	Topic     string
	Partition int
	Host      string
	Network   string
}

type MassageRead struct {
	ID      string `json:"ID"`
	Massage string `json:"Massage"`
}

type MessageWrite struct {
	ID     string
	Status models.Status
}

func main() {
	time.Sleep(20 * time.Second) // give kafka time to start
	if err := godotenv.Overload("env/.env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	confWrite := Config{
		Topic:     "topic2",
		Partition: 0,
		Host:      os.Getenv("KAFKA_ADDRESS"),
		Network:   "tcp",
	}
	writer, err := Connection(confWrite)
	if err != nil {
		fmt.Println("Error connecting to kafka: %v", err)
	}

	confRead := kafka.ReaderConfig{
		Topic:     "topic1",
		Partition: 0,
		GroupID:   "g1",
		Brokers:   []string{os.Getenv("KAFKA_ADDRESS")},
	}

	reader := kafka.NewReader(confRead)

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("Error reading message: %s\n", err)
			break
		}
		fmt.Println("Message: ", string(msg.Value))

		r := MassageRead{}
		err = json.Unmarshal(msg.Value, &r)
		if err != nil {
			fmt.Printf("Error decoding message: %s\n", err)
		}

		w := MessageWrite{
			ID:     r.ID,
			Status: models.StatusCompleted,
		}

		value, err := json.Marshal(w)
		if err != nil {
			fmt.Println("failed to serialize structure: %v", err)
		}

		_, err = writer.WriteMessages(
			kafka.Message{Value: value},
		)
		if err != nil {
			fmt.Println("failed to write messages: %v", err)
		}
	}
}

func Connection(conf Config) (*kafka.Conn, error) {
	conn, err := kafka.DialLeader(context.Background(), conf.Network, conf.Host, conf.Topic, conf.Partition)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
