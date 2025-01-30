package main

import (
	"context"
	"gate/api"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/segmentio/kafka-go"
)

type Server struct {
	kafkaReaderAck *kafka.Reader
	kafkaWriter    *kafka.Writer
}

func NewServer() api.ServerInterface {

	return &Server{
		kafkaReaderAck: getKafkaReaderAck(),
		kafkaWriter:    getKafkaWriter(),
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

	err := sendPrintJobToKafka(s.kafkaWriter, printJob)
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

func (s *Server) GetPrintStatus(ctx echo.Context) error {
	log.Println("Received request to /print/status")

	confirmation, err := readPrintJobAckFromKafka(s.kafkaReaderAck)
	if err != nil {
		log.Print(err)
		return ctx.String(http.StatusInternalServerError, "Failed to read print job ack from kafka")
	}
	if confirmation.Status == "" {
		msg := "No print jobs completed"
		return ctx.JSON(http.StatusOK, api.PrintStatus{
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

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:20000", "http://example.com"},                     // Allowed origins
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete}, // Allowed methods
	}))

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
	// Set a timeout to shutdown gracefully to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	closeKafkaReader(server.(*Server).kafkaReaderAck)
	closeKafkaWriter(server.(*Server).kafkaWriter)

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}
	log.Println("Server shutdown complete.")
}
