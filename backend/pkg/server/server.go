package server

import (
	"backend/config"
	"backend/pb"
	"backend/pkg/cache"
	"backend/pkg/db"
	"backend/pkg/db/model"
	"backend/pkg/repo"
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
	"os/signal"
	"strconv"
	"syscall"
)

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {

	return &Server{
		config: config,
	}
}

func initTerminationChan(proxyService *service.ProxyService) {
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-termChan
		log.Print("Received termination signal, storing cache...")
		proxyService.StoreCache()
		os.Exit(1)

	}()
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

	cacheCap, err := strconv.Atoi(server.config.CacheCap)
	utils.Handle(err)

	lruCache := cache.NewLRUCache(cacheCap)
	gormDB, err := db.Init(server.config.DB)

	if err != nil {
		log.Fatal("Could not connect to the database", err.Error())
		return
	}

	proxyService := &service.ProxyService{
		EthClient:   ethClient,
		Cache:       lruCache,
		Config:      server.config,
		ProxyNftAbi: nftAbi,
		HolderRepo: repo.Repo[model.HolderEntity]{
			DB: gormDB,
		},
	}

	err = proxyService.LoadCacheFromDB()

	if err != nil {
		log.Fatal("Could not load cache from the database ", err.Error())
		return
	}

	go proxyService.StartCacheValidation()
	initTerminationChan(proxyService)

	pb.RegisterEthProxyServer(s, proxyService)

	log.Print("Server staring on port: ", server.config.Port)
	if err := s.Serve(listener); err != nil {
		log.Fatal("Failed to serve")
	}
}
