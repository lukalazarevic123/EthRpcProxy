package service

import (
	"backend/pb"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

type ProxyService struct {
	pb.UnimplementedEthProxyServer
	EthClient *ethclient.Client
}

type SendTransactionResponse struct {
	TxnHash string `json:"txnHash"`
}

func (ps *ProxyService) getLatestBlockNumber() uint64 {
	log.Print("pokusaj")
	blockNumber, err := ps.EthClient.BlockNumber(context.Background())
	if err != nil {
		log.Println("Error fetching latest block header:", err)
		return 0
	}
	log.Print(blockNumber)
	return blockNumber
}

func (ps *ProxyService) EthSendTransaction(ctx context.Context, args *pb.SendTransactionRequest) (*pb.TransactionReceipt, error) {

	blockNum := ps.getLatestBlockNumber()
	log.Print(blockNum)
	return &pb.TransactionReceipt{
		Hash: fmt.Sprint(blockNum),
	}, nil
}
