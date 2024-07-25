package server

import (
	"backend/config"
	"backend/pb"
	"backend/pkg/cache"
	"backend/pkg/service"
	"backend/pkg/utils"
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"strconv"
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

	abiData, err := os.ReadFile("./abi/ProxyNFT.json")
	utils.Handle(err)

	nftAbi, err := abi.JSON(bytes.NewReader(abiData))
	utils.Handle(err)
	log.Print(nftAbi)
	cacheCap, err := strconv.Atoi(server.config.CacheCap)
	utils.Handle(err)

	lruCache := cache.NewLRUCache(cacheCap)

	proxyService := &service.ProxyService{
		EthClient:   ethClient,
		Cache:       lruCache,
		Config:      server.config,
		ProxyNftAbi: nftAbi,
	}
	go proxyService.StartCacheValidation()
	pb.RegisterEthProxyServer(s, proxyService)

	log.Print("Server staring on port: ", server.config.Port)
	if err := s.Serve(listener); err != nil {
		log.Fatal("Failed to serve")
	}
}
