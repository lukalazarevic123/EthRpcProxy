package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	Port            string
	Interval        string
	EthRpcUrl       string
	CacheCap        string
	ProxyNftAddress string
	DB              *DBConfig
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
	SslMode  string
}

func NewConfig() *Config {
	port, err := strconv.Atoi(readEnvVar("DB_PORT"))

	if err != nil {
		return nil
	}

	return &Config{
		Port:            readEnvVar("PORT"),
		Interval:        readEnvVar("INTERVAL"),
		EthRpcUrl:       readEnvVar("ETH_RPC_URL"),
		CacheCap:        readEnvVar("CACHE_CAP"),
		ProxyNftAddress: readEnvVar("PROXY_NFT_ADDRESS"),
		DB: &DBConfig{
			User:     readEnvVar("DB_USER"),
			Dbname:   readEnvVar("DB_NAME"),
			Host:     readEnvVar("DB_HOST"),
			SslMode:  readEnvVar("DB_SSL"),
			Password: readEnvVar("DB_PASSWORD"),
			Port:     port,
		},
	}
}

func readEnvVar(name string) string {
	godotenv.Load(".env")

	return os.Getenv(name)
}
