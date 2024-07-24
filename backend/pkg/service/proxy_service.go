package service

import (
	"backend/pb"
	"backend/pkg/utils"
	"context"
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
	blockNumber, err := ps.EthClient.BlockNumber(context.Background())
	if err != nil {
		log.Println("Error fetching latest block header:", err)
		return 0
	}
	return blockNumber
}

func (ps *ProxyService) sendTransaction(args *pb.SendTransactionRequest) (string, error) {
	var txnHash string
	err := ps.EthClient.Client().Call(&txnHash, "eth_sendTransaction", args)

	return txnHash, err
}

func (ps *ProxyService) EthSendTransaction(ctx context.Context, args *pb.SendTransactionRequest) (*pb.TransactionReceipt, error) {
	txnHash, err := ps.sendTransaction(args)

	utils.Handle(err)
	return &pb.TransactionReceipt{
		Hash: txnHash,
	}, nil
}
