package util

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	AuthServerPort        string `mapstructure:"AUTH_SERVER_PORT"`
	EnableAPM             int16  `mapstructure:"ENABLE_APM"`
	NodeEnv               string `mapstructure:"NODE_ENV"`
	SaltHash              int16  `mapstructure:"SALT_HASH"`
	ClientUrl             string `mapstructure:"CLIENT_URL"`
	CloudName             string `mapstructure:"CLOUD_NAME"`
	CloudAPIKey           string `mapstructure:"CLOUD_API_KEY"`
	CloudAPISecret        string `mapstructure:"CLOUD_API_SECRET"`
	APIGatewayUrl         string `mapstructure:"API_GATEWAY_URL"`
	ElasticSearchUrl      string `mapstructure:"ELASTIC_SEARCH_URL"`
	ElasticAPMServerUrl   string `mapstructure:"ELASTIC_APM_SERVER_URL"`
	RabbitMQEndPoint      string `mapstructure:"RABBITMQ_ENDPOINT"`
	DbSource              string `mapstructure:"DB_SOURCE"`
	LogEnable             bool   `mapstructure:"LOG_ENABLE"`
	MaxConnect            int    `mapstructure:"MAX_CONNECT"`
	IdleConnect           int    `mapstructure:"IDLE_CONNECT"`
	MaxLifeTime           int    `mapstructure:"MAX_LIFE_TIME"`
	ElasticAPMSecretToken string `mapstructure:"ELASTIC_APM_SECRET_TOKEN"`
	EmailExchangeName     string `mapstructure:"EMAIL_EXCHANGE_NAME"`
	EmailQueueName        string `mapstructure:"EMAIL_QUEUE_NAME"`
	EmailRoutingKey       string `mapstructure:"EMAIL_ROUTING_KEY"`
	GatewayToken          string `mapstructure:"GATEWAY_JWT_TOKEN"`
	JWTSecret             string `mapstructure:"JWT_SECRET"`
	BasePath              string `mapstructure:"BASE_PATH"`
	Auth                  string `mapstructure:"AUTH"`
}

func LoadConfig(path string) (config Config, err error) {

	fmt.Println("load func called")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	err = viper.Unmarshal(&config)
	return
}
