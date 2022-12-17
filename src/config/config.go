package config

import (
	"fmt"
	"log"
	"sync"

	"gopkg.in/yaml.v3"
)

type Trigger struct {
	Event   string
	Command string
}

type QueueConfig struct {
	Url    string
	Region string
	KeyId  string
	Secret string
}

type Config struct {
	Queue    QueueConfig
	Triggers []Trigger
}

type ConfigData struct {
	Queue struct {
		Url    string
		Region string
		KeyId  string `yaml:"keyId"`
		Secret string
	}

	Triggers []map[string]string
}

var lock = &sync.Mutex{}
var reader ConfigReader = &FileReader{}
var instance = &Config{}

const CONFIG = "/etc/cloud-tunnel.yaml"

/*
 * Load a configuration file and initialize the Config singleton
 */
func Load(file ...string) error {
	lock.Lock()
	defer lock.Unlock()

	var path = CONFIG
	if len(file) > 0 && file[0] != "" {
		path = file[0]
	}

	log.Println("Loading application configuration:", path)

	// Read the config file yaml
	yaml, err := reader.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Error loading configuration: %s", err)
	}

	// Populate the config object
	config, err := parse(yaml)
	if err != nil {
		return fmt.Errorf("Error parsing configuration: %s", err)
	}

	instance.Queue = config.Queue
	instance.Triggers = config.Triggers

	return nil
}

/**
* Parse the yaml string and map it into the configuration object
 */
func parse(content []byte) (*Config, error) {
	var data = &ConfigData{}
	err := yaml.Unmarshal(content, data)
	if err != nil {
		return nil, err
	}

	// Populate the config object
	var config = &Config{}
	config.Queue = QueueConfig(data.Queue)
	config.Triggers = make([]Trigger, 0)

	for _, trigger := range data.Triggers {
		for event, command := range trigger {
			config.Triggers = append(config.Triggers, Trigger{
				Event:   event,
				Command: command,
			})
		}
	}

	return config, nil
}

/*
 * Return the configuration instance
 */
func Instance() *Config {
	return instance
}
