package service

import (
	"backend/pkg/requests"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

type ProxyService struct {
	EthClient *ethclient.Client
}

type SendTransactionResponse struct {
	TxnHash string `json:"txnHash"`
}

func (ps *ProxyService) EthSendTransaction(args requests.SendTransactionArgs, reply *SendTransactionResponse) error {
	// Replace with actual implementation logic
	*reply = SendTransactionResponse{
		TxnHash: "asdasdasdasdasd",
	}
	log.Print(args)
	return nil
}
