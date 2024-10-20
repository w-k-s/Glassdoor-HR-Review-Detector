package server

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		ListenAddress string `mapstructure:"listenAddress"`
	} `mapstructure:"server"`
	Database struct {
		Url string `mapstructure:"url"`
	} `mapstructure:"database"`
	S3 struct {
		Bucket string `mapstructure:"bucket"`
		Region string `mapstructure:"region"`
	} `mapstructure:"s3"`
	OpenAi struct {
		ApiKey string `mapstructure:"apiKey"`
	} `mapstructure:"openai"`
	Migrations struct {
		Directory string `mapstructure:"directory"`
	} `mapstructure:"migrations"`
	Feedback struct {
		Upload struct {
			Frequency struct {
				Hours int `mapstructure:"hours"`
			} `mapstructure:"frequency"`
		} `mapstructure:"upload"`
	} `mapstructure:"feedback"`
	Inference struct {
		Api struct {
			Endpoint string `mapstructure:"endpoint"`
		} `mapstructure:"api"`
	} `mapstructure:"inference"`
}

func ReadConfig() *Config {
	applicationDirectory := applicationDirectory()

	viper.SetDefault("server.listenAddress", ":3000")
	viper.SetDefault("database.url", "file::memory:?cache=shared")
	viper.SetDefault("s3.bucket", "glassdoor-hr-review-detector")
	viper.SetDefault("s3.region", "ap-south-1")
	viper.SetDefault("migrations.directory", filepath.Join(applicationDirectory, "migrations"))
	viper.SetDefault("feedback.upload.frequency.hours", 1)
	viper.SetDefault("inference.api.endpoint", "http://localhost:8501/v1/models/glassdoor_hr_review_detector:predict")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(applicationDirectory)
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No config file found. Using defaults and environment variables.")
		} else {
			log.Fatalf("Error reading config file: %s", err)
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("GDHR")
	viper.BindEnv("server.listenAddress", "GDHR_LISTEN_ADDRESS")
	viper.BindEnv("database.url", "GDHR_DB_URL")
	viper.BindEnv("s3.bucket", "GDHR_S3_BUCKET")
	viper.BindEnv("s3.region", "AWS_REGION")
	viper.BindEnv("openai.apiKey", "OPENAI_API_KEY")
	viper.BindEnv("migrations.directory", "GDHR_MIGRATIONS_DIRECTORY")
	viper.BindEnv("feedback.upload.frequency.hours", "GDHR_FEEDBACK_UPLOAD_FREQUENCY_HOURS")
	viper.BindEnv("inference.api.endpoint", "GDHR_INFERENCE_API_ENDPOINT")

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	fmt.Printf("Server Listen Address: %s\n", config.Server.ListenAddress)
	fmt.Printf("Database URL: %s\n", config.Database.Url)
	fmt.Printf("S3 Bucket: %s\n", config.S3.Bucket)
	fmt.Printf("S3 Region: %s\n", config.S3.Region)
	fmt.Printf("OpenAI API Key: %s\n", config.OpenAi.ApiKey)
	fmt.Printf("Migrations Directory: %s\n", config.Migrations.Directory)
	fmt.Printf("Feedback Upload Frequency: Every %d Hour(s)\n", config.Feedback.Upload.Frequency.Hours)
	fmt.Printf("Inference API Endpoint: %s\n", config.Inference.Api.Endpoint)

	return &config
}

func applicationDirectory() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting user home directory: %v", err)
	}
	return filepath.Join(homeDir, ".glassdoor-hr-review-detector")
}
