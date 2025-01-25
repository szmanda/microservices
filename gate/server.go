package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
	"gate/api"

	"github.com/segmentio/kafka-go"
	"github.com/labstack/echo/v4"
)

type Server struct {
}

func NewServer() api.ServerInterface{
    return &Server{}
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

func StartServer() {
    server := NewServer()
	e := echo.New()
	api.RegisterHandlers(e, server)
	
    log.Printf("Service is running on port %s...\n", appConfig.Port)
	if err := e.Start(":"+appConfig.Port); err != nil && err != http.ErrServerClosed {
 	   log.Fatalf("Could not listen on :%s %v\n", appConfig.Port, err)
    }
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