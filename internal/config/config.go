package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port     string   `mapstructure:"port"`
	LogLevel string   `mapstructure:"log_level"`
	Database Database `mapstructure:"database"`
	Kafka    Kafka    `mapstructure:"kafka"`
	External External `mapstructure:"external"`
	Payment  Payment  `mapstructure:"payment"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type Kafka struct {
	Brokers []string `mapstructure:"brokers"`
	Topics  Topics   `mapstructure:"topics"`
}

type Topics struct {
	BookingEvents     string `mapstructure:"booking_events"`
	PaymentEvents     string `mapstructure:"payment_events"`
	SeatSelectEvents  string `mapstructure:"seat_select_events"`
}

type External struct {
	HackloadBaseURL string `mapstructure:"hackload_base_url"`
}

type Payment struct {
	GatewayURL string `mapstructure:"gateway_url"`
	TeamSlug   string `mapstructure:"team_slug"`
	Password   string `mapstructure:"password"`
}

func Load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("/app")  // для Docker контейнера

	// Устанавливаем значения по умолчанию
	viper.SetDefault("port", "8081")
	viper.SetDefault("log_level", "info")
	viper.SetDefault("database.host", "biletter-postgres")  // имя сервиса в docker-compose
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "biletter_user")
	viper.SetDefault("database.password", "biletter_pass")
	viper.SetDefault("database.dbname", "biletter_db")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("kafka.brokers", []string{"biletter-kafka:29092"})
	viper.SetDefault("external.hackload_base_url", "http://localhost:8080")
	viper.SetDefault("payment.gateway_url", "https://hub.hackload.kz/payment-provider/common/api/v1")
	viper.SetDefault("payment.team_slug", "metaload-akbori")
	
	// Привязываем переменные окружения
	viper.BindEnv("port", "PORT")
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.user", "DB_USERNAME")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("payment.password", "PAYMENT_PASSWORD")
	viper.BindEnv("external.hackload_base_url", "HACKLOAD_BASE_URL")

	// Пытаемся прочитать конфигурационный файл (опционально)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using defaults and environment variables")
		} else {
			log.Printf("Error reading config file: %v", err)
		}
	} else {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode config: %v", err)
	}

	return &config
}