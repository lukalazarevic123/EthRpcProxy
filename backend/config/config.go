package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port      string
	Interval  string
	EthRpcUrl string
}

func NewConfig() *Config {
	return &Config{
		Port:      readEnvVar("PORT"),
		Interval:  readEnvVar("INTERVAL"),
		EthRpcUrl: readEnvVar("ETH_RPC_URL"),
	}
}

func readEnvVar(name string) string {
	godotenv.Load(".env")

	return os.Getenv(name)
}
