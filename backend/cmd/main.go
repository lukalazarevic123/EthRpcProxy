package main

import (
	"backend/config"
	"backend/pkg/server"
)

func main() {
	cfg := config.NewConfig()
	srv := server.NewServer(cfg)

	srv.Start()
}
