package main

import (
	"context"
	"encoding/json"
	"gate/api"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/segmentio/kafka-go"
)

type Server struct {
	kafkaReaderAck *kafka.Reader
}

func NewServer() api.ServerInterface {

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{appConfig.KafkaBroker},
		Topic:       appConfig.KafkaTopicAck,
		GroupID:     appConfig.ServiceName,
		Partition:   0,
		MinBytes:    10e3,
		MaxBytes:    10e6,
		StartOffset: 0,
	})

	return &Server{
		kafkaReaderAck: r,
	}
}

func (s *Server) GetHello(ctx echo.Context) error {
	log.Println("Received request to /hello")

	msg := "Hello from the gate service!"
	response := api.Response{
		Message: &msg,
		Service: &appConfig.ServiceName,
	}

	return ctx.JSON(http.StatusOK, response)
}

func (s *Server) PostPrint(ctx echo.Context) error {
	log.Println("Received request to /print")

	var printJob api.PrintJob
	if err := ctx.Bind(&printJob); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		return ctx.String(http.StatusBadRequest, "Invalid request body")
	}
	log.Printf("Received print job: %+v\n", printJob)

	err := SendPrintJobToKafka(printJob)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to send print job to queue")
	}
	msg := "Print job sent to queue"
	response := api.Response{
		Message: &msg,
		Service: &appConfig.ServiceName,
	}

	return ctx.JSON(http.StatusOK, response)

}

// Kafka confirmation message
type PrintConfirmation struct {
	DocumentName string    `json:"document_name"`
	Status       string    `json:"status"`
	Timestamp    time.Time `json:"timestamp"`
}

func (s *Server) GetPrintStatus(ctx echo.Context) error {
	log.Println("Received request to /print/status")

	msgCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	// msgCtx := context.Background() // no timeout -- a blocking read

	m, err := s.kafkaReaderAck.ReadMessage(msgCtx)
	if err != nil {
		log.Printf("No message in the topic: %v\n", err)
		msg := "print not complete"
		return ctx.JSON(http.StatusOK, api.PrintStatus{
			Message: &msg,
			Service: &appConfig.ServiceName,
		})
	}

	log.Printf("Message at offset %d: %s\n", m.Offset, string(m.Value))

	var confirmation PrintConfirmation
	if err := json.Unmarshal(m.Value, &confirmation); err != nil {
		log.Printf("Error unmarshalling confirmation message: %v\n", err)
		msg := "error processing ack message"
		return ctx.JSON(http.StatusInternalServerError, api.PrintStatus{
			Message: &msg,
			Service: &appConfig.ServiceName,
		})
	}

	msg := "Print Job for '" + confirmation.DocumentName + "' is " + confirmation.Status + " at " + confirmation.Timestamp.UTC().String()
	return ctx.JSON(http.StatusOK, api.PrintStatus{
		Message: &msg,
		Service: &appConfig.ServiceName,
	})
}

func (s *Server) PostNipChecker(ctx echo.Context) error {
	log.Println("Received request to /nip_checker")
	var nipRequest api.NipRequest
	if err := ctx.Bind(&nipRequest); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		return ctx.String(http.StatusBadRequest, "Invalid request body")
	}

	nip := nipRequest.Nip
	nipInfo, err := getNipData(nip)
	if err != nil {
		log.Printf("Error getting NIP data: %v", err)
		return ctx.String(http.StatusInternalServerError, "Failed to fetch NIP data")
	}
	return ctx.JSON(http.StatusOK, nipInfo)

}

func StartServer() {
	server := NewServer()
	e := echo.New()
	api.RegisterHandlers(e, server)

	// Start the server in a goroutine so that it doesn't block the main thread.
	go func() {
		log.Printf("Service is running on port %s...\n", appConfig.Port)
		if err := e.Start(":" + appConfig.Port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on :%s %v\n", appConfig.Port, err)
		}
	}()

	quit := make(chan os.Signal, 1) // Wait for interrupt signal
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Close the Kafka reader before exiting
	if server.(*Server).kafkaReaderAck != nil {
		if err := server.(*Server).kafkaReaderAck.Close(); err != nil {
			log.Printf("Error closing kafka reader: %v", err)
		}
	}
	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}
	log.Println("Server shutdown complete.")

}

func SendPrintJobToKafka(printJob api.PrintJob) error {

	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{appConfig.KafkaBroker},
		Topic:   appConfig.KafkaTopic,
	})
	defer kafkaWriter.Close()

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
