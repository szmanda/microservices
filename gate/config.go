package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServiceName   string `yaml:"servcie_name"`
	Port          string `yaml:"port"`
	KafkaBroker   string `yaml:"kafka_broker"`
	KafkaTopic    string `yaml:"kafka_topic"`
	KafkaTopicAck string `yaml:"kafka_topic_ack"`
}

var appConfig Config

func LoadConfig(filename string) error {
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
