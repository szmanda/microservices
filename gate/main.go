package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ServiceName                string        `yaml:"servcie_name"`
	Port                       string        `yaml:"port"`
    KafkaBroker string `yaml:"kafka_broker"`
	KafkaTopic  string `yaml:"kafka_topic"`
}

var appConfig Config

type Response struct {
	Message string `json:"message"`
	Service string `json:"service"`
}

type PrintJob struct {
	DocumentName  string   `json:"document_name"`
	PaperSize     string   `json:"paper_size"`
	Orientation   string   `json:"orientation"`
	Copies        int      `json:"copies"`
	Printers      []string `json:"printers"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to /hello")

	response := Response{
		Message: "Hello world!",
		Service: appConfig.ServiceName,
	}

    w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding JSON: %v", err)
        return
    }
}

func printHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to /print")

	// Example print job from request body
	var printJob PrintJob
	err := json.NewDecoder(r.Body).Decode(&printJob)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		log.Printf("Error decoding JSON: %v", err)
		return
	}
	log.Printf("Received print job: %+v\n", printJob)

	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{appConfig.KafkaBroker},
		Topic:   appConfig.KafkaTopic,
	})
	defer kafkaWriter.Close()

	message, err := json.Marshal(printJob)
	if err != nil {
		http.Error(w, "Failed to encode print job", http.StatusInternalServerError)
		log.Printf("Error encoding print job: %v", err)
		return
	}
	
	err = kafkaWriter.WriteMessages(context.Background(), kafka.Message{
		Value: message,
	})

	if err != nil {
		http.Error(w, "Failed to write message to kafka", http.StatusInternalServerError)
		log.Printf("Error sending message to kafka: %v", err)
		return
	}


	response := Response{
		Message: "Print job sent to queue",
		Service: appConfig.ServiceName,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding JSON: %v", err)
		return
	}
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

	err := loadConfig("gate/parameters.yml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	log.Printf("Loaded config: %+v\n", appConfig)


	r := mux.NewRouter()
	r.HandleFunc("/hello", helloHandler).Methods("GET")
    r.HandleFunc("/print", printHandler).Methods("POST")

	server := &http.Server{
		Addr:    ":80",
		Handler: r,
	}
	log.Println("Service is running on port 80...")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on :80 %v\n", err)
	}
}