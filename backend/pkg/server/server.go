package server

import (
	"backend/config"
	"backend/pb"
	"backend/pkg/service"
	"backend/pkg/utils"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {

	return &Server{
		config: config,
	}
}

func (server *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	utils.Handle(err)

	s := grpc.NewServer()
	reflection.Register(s)

	ethClient, err := ethclient.Dial(server.config.EthRpcUrl)
	utils.Handle(err)

	pb.RegisterEthProxyServer(s, &service.ProxyService{
		EthClient: ethClient,
	})

	log.Print("Server staring on port: ", server.config.Port)
	if err := s.Serve(listener); err != nil {
		log.Fatal("Failed to serve")
	}
}
