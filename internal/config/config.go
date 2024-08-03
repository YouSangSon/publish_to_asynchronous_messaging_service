package config

import (
	"encoding/json"
	"os"

	"go.uber.org/zap"
)

var (
	ModeConfig Configuration
)

const (
	envConfigPath   = "configJson/config.json"
	envPubSubPath   = "configJson/pub_sub.json"
	envKafkaPath    = "configJson/kafka.json"
	envRabbitmqPath = "configJson/rabbitmq.json"
)

type Configurations struct {
	Dev Configuration `json:"dev"`
	Pro Configuration `json:"prod"`
}

type Configuration struct {
	Logging    zap.Config      `json:"logging"`
	PubSubs    *pubSub         `json:"pub_subs,omitempty"`
	KafkaStore *storeConfig    `json:"kafka_store,omitempty"`
	RabbitMq   *rabbitMqConfig `json:"rabbitmq,omitempty"`
}

type kafkaConfig struct {
	Dev kafkaStoreConfig `json:"dev"`
	Pro kafkaStoreConfig `json:"prod"`
}

type kafkaStoreConfig struct {
	KafkaStore *storeConfig `json:"kafka_store"`
}

type storeConfig struct {
	Host  string `json:"host"`
	Port  int    `json:"port"`
	Topic string `json:"topic"`
}

type pubSubsConfig struct {
	Dev pubSubs `json:"dev"`
	Pro pubSubs `json:"prod"`
}

type pubSubs struct {
	PubSub *pubSub `json:"pub_subs"`
}

type pubSub struct {
	Project string `json:"project"`
	Topic   string `json:"topic"`
	Subname string `json:"subname"`
}

type rabbitMqConfig struct {
	Server   string `json:"server"`
	Vhost    string `json:"vhost"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// InitConfig initializes the configuration
func InitConfig(mode, runMode string) error {
	var configurations Configurations

	if err := loadConfig(envConfigPath, &configurations); err != nil {
		return err
	}

	switch runMode {
	case "pubsub":
		var pubSubsConfig pubSubsConfig

		if err := loadConfig(envPubSubPath, &pubSubsConfig); err != nil {
			return err
		}

		configurations.Pro.PubSubs = pubSubsConfig.Pro.PubSub
		configurations.Dev.PubSubs = pubSubsConfig.Dev.PubSub
	case "kafka":
		var kafkaConfig kafkaConfig

		if err := loadConfig(envKafkaPath, &kafkaConfig); err != nil {
			return err
		}

		configurations.Pro.KafkaStore = kafkaConfig.Pro.KafkaStore
		configurations.Dev.KafkaStore = kafkaConfig.Dev.KafkaStore
	case "rabbitmq":
		var rabbitMqConfig rabbitMqConfig

		if err := loadConfig(envRabbitmqPath, &rabbitMqConfig); err != nil {
			return err
		}

		configurations.Pro.RabbitMq = &rabbitMqConfig
		configurations.Dev.RabbitMq = &rabbitMqConfig
	default:
	}

	switch mode {
	case "prod":
		ModeConfig = configurations.Pro
	case "dev":
		ModeConfig = configurations.Dev
	default:
		ModeConfig = configurations.Pro
	}

	return nil
}

// loadConfig loads the configuration from the file
func loadConfig(configPath string, configData interface{}) error {
	configFile, err := os.Open(configPath)
	if err != nil {
		return err
	}

	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&configData); err != nil {
		return err
	}

	configFile.Close()

	return nil
}
