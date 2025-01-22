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
	ServiceName string `yaml:"servcie_name"`
	KafkaBroker string `yaml:"kafka_broker"`
	KafkaTopic  string `yaml:"kafka_topic"`
}

var appConfig Config

type PrintJob struct {
	DocumentName  string   `json:"document_name"`
	PaperSize     string   `json:"paper_size"`
	Orientation   string   `json:"orientation"`
	Copies        int      `json:"copies"`
	Printers      []string `json:"printers"`
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
	err := loadConfig("service/parameters.yml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	log.Printf("Loaded config: %+v\n", appConfig)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{appConfig.KafkaBroker},
		Topic:     appConfig.KafkaTopic,
		GroupID:   appConfig.ServiceName, // Assign a group ID
		Partition: 0,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	defer r.Close()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue // Skip and continue with next message
		}
		log.Printf("Message at offset %d: %s\n", m.Offset, string(m.Value))

		var job PrintJob
		err = json.Unmarshal(m.Value, &job)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue // Skip and continue with next message
		}

		log.Printf("Received print job: %+v\n", job)
		log.Println("Printing...")

		time.Sleep(5 * time.Second)
		log.Println("Print job completed.")
		log.Printf("Message processing finished for offset %d\n", m.Offset)
	}
}