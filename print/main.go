package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ServiceName   string `yaml:"servcie_name"`
	KafkaBroker   string `yaml:"kafka_broker"`
	KafkaTopic    string `yaml:"kafka_topic"`
	KafkaTopicAck string `yaml:"kafka_topic_ack"` // Added Ack Topic
}

var appConfig Config

type PrintJob struct {
	DocumentName string   `json:"document_name"`
	PaperSize    string   `json:"paper_size"`
	Orientation  string   `json:"orientation"`
	Copies       int      `json:"copies"`
	Printers     []string `json:"printers"`
}

type PrintConfirmation struct {
	DocumentName string    `json:"document_name"`
	Status       string    `json:"status"`
	Timestamp    time.Time `json:"timestamp"`
}

func loadConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening config file: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&appConfig); err != nil {
		return fmt.Errorf("error decoding config file: %w", err)
	}

	return nil
}

func main() {
	err := loadConfig("parameters.yml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	log.Printf("Loaded config: %+v\n", appConfig)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{appConfig.KafkaBroker},
		Topic:       appConfig.KafkaTopic,
		GroupID:     appConfig.ServiceName,
		Partition:   0,
		MinBytes:    10e3,
		MaxBytes:    10e6,
		StartOffset: 0,
	})
	defer r.Close()

	w := &kafka.Writer{
		Addr:     kafka.TCP(appConfig.KafkaBroker),
		Topic:    appConfig.KafkaTopicAck,
		Balancer: &kafka.LeastBytes{},
	}
	defer w.Close()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}
		log.Printf("Message at offset %d: %s\n", m.Offset, string(m.Value))

		var job PrintJob
		err = json.Unmarshal(m.Value, &job)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		log.Printf("Received print job: %+v\n", job)
		log.Println("Printing...")

		time.Sleep(5 * time.Second)
		log.Println("Print job completed.")

		confirmation := PrintConfirmation{
			DocumentName: job.DocumentName,
			Status:       "success",
			Timestamp:    time.Now(),
		}

		confirmationBytes, err := json.Marshal(confirmation)
		if err != nil {
			log.Printf("Error marshalling confirmation message: %v", err)
			continue
		}

		err = w.WriteMessages(context.Background(), kafka.Message{
			Value: confirmationBytes,
		})

		if err != nil {
			log.Printf("Error sending confirmation message: %v", err)
			continue
		}

		log.Printf("Confirmation message sent to topic %s: %s\n", appConfig.KafkaTopicAck, string(confirmationBytes))
	}
}
