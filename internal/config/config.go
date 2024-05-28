package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const (
	defaultGRPCPort = "443"

	defaultHTTPPort               = "8000"
	defaultHTTPRWTimeout          = 10 * time.Second
	defaultHTTPMaxHeaderMegabytes = 1

	envLocal  = "local"
	envDev    = "dev"
	envProd   = "prod"
	authority = "notification-microservice"
)

type (
	Config struct {
		Environment  string
		GRPC         GRPCConfig `mapstructure:"grpc"`
		SMTP         SMTPConfig `mapstructure:"smtp"`
		Authority    string
		QRs          MicroserviceConfig `mapstructure:"qrMicroservice"`
		Users        MicroserviceConfig `mapstructure:"userMicroservice"`
		Reservations MicroserviceConfig `mapstructure:"reservationMicroservice"`
	}

	SMTPConfig struct {
		Host     string
		Port     int
		Username string
		Password string
		Sender   string
	}

	MicroserviceConfig struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	}

	GRPCConfig struct {
		Host    string        `mapstructure:"host"`
		Port    string        `mapstructure:"port"`
		Timeout time.Duration `mapstructure:"timeout"`
	}
)

func Init(configsDir, envDir string) (*Config, error) {
	populateDefaults()
	loadEnvVariables(envDir)
	if err := parseConfigFile(configsDir, ""); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	setFromEnv(&cfg)
	return &cfg, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("grpc", &cfg.GRPC); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("qrMicroservice", &cfg.QRs); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("reservationMicroservice", &cfg.Reservations); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("userMicroservice", &cfg.Users); err != nil {
		return err
	}

	return viper.UnmarshalKey("smtp", &cfg.SMTP)
}

func setFromEnv(cfg *Config) {
	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	cfg.SMTP.Port, _ = strconv.Atoi(os.Getenv("SMTP_PORT"))
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")
	cfg.SMTP.Sender = os.Getenv("SMTP_SENDER")

	cfg.Authority = authority
	cfg.GRPC.Host = os.Getenv("GRPC_HOST")

	cfg.Environment = envDev
}

func parseConfigFile(folder, env string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	viper.SetConfigName(env)

	return viper.MergeInConfig()
}

func loadEnvVariables(envPath string) {
	err := godotenv.Load(envPath)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

}

func populateDefaults() {
	viper.SetDefault("grpc.port", defaultGRPCPort)
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("http.max_header_megabytes", defaultHTTPMaxHeaderMegabytes)
	viper.SetDefault("http.timeouts.read", defaultHTTPRWTimeout)
	viper.SetDefault("http.timeouts.write", defaultHTTPRWTimeout)
}
