package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gate/api"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func getKafkaReaderAck() *kafka.Reader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{appConfig.KafkaBroker},
		Topic:       appConfig.KafkaTopicAck,
		GroupID:     appConfig.ServiceName,
		Partition:   0,
		MinBytes:    10e3,
		MaxBytes:    10e6,
		StartOffset: 0,
	})
	return r
}

func getKafkaWriter() *kafka.Writer {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{appConfig.KafkaBroker},
		Topic:   appConfig.KafkaTopic,
	})
	return w
}

func sendPrintJobToKafka(kafkaWriter *kafka.Writer, printJob api.PrintJob) error {
	message, err := json.Marshal(printJob)
	if err != nil {
		log.Printf("Error encoding print job: %v", err)
		return err
	}

	err = kafkaWriter.WriteMessages(context.Background(), kafka.Message{
		Value: message,
	})

	if err != nil {
		log.Printf("Error sending message to kafka: %v", err)
		return err
	}
	return nil
}

// Kafka confirmation message
type PrintConfirmation struct {
	DocumentName string    `json:"document_name"`
	Status       string    `json:"status"`
	Timestamp    time.Time `json:"timestamp"`
}

func readPrintJobAckFromKafka(kafkaReader *kafka.Reader) (PrintConfirmation, error) {

	msgCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	// msgCtx := context.Background() // no timeout -- a blocking read

	var confirm PrintConfirmation
	m, err := kafkaReader.ReadMessage(msgCtx)
	if err != nil {
		return confirm, nil // no new message
	}
	log.Printf("Message at offset %d: %s\n", m.Offset, string(m.Value))

	if err := json.Unmarshal(m.Value, &confirm); err != nil {
		return confirm, fmt.Errorf("error unmarshalling JSON from Kafka message: %w, message: %s", err, string(m.Value))
	}

	return confirm, nil
}

func closeKafkaReader(reader *kafka.Reader) {
	if reader == nil {
		return
	}
	err := reader.Close()
	if err != nil {
		log.Printf("Error closing kafka reader: %v", err)
	}
}

func closeKafkaWriter(writer *kafka.Writer) {
	if writer == nil {
		return
	}
	err := writer.Close()
	if err != nil {
		log.Printf("Error closing kafka writer: %v", err)
	}
}
