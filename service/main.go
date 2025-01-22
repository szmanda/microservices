package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ServiceName                string        `yaml:"servcie_name"`
	Port                       string        `yaml:"port"`
}
var appConfig Config

type Response struct {
	Message string `json:"message"`
	Service string `json:"service"`
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

	r := mux.NewRouter()
	r.HandleFunc("/hello", helloHandler).Methods("GET")

	server := &http.Server{
		Addr:    ":80",
		Handler: r,
	}
	log.Println("Service is running on port 80...")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on :80 %v\n", err)
	}
}
