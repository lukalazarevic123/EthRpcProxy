package service

import (
	"backend/config"
	"backend/pb"
	"backend/pkg/cache"
	"backend/pkg/utils"
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"strconv"
	"time"
)

type ProxyService struct {
	pb.UnimplementedEthProxyServer
	EthClient *ethclient.Client
	Cache     *cache.LRUCache
	Config    *config.Config
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

func (ps *ProxyService) StartCacheValidation() {
	intervalNum, err := strconv.Atoi(ps.Config.Interval)
	utils.Handle(err)

	ticker := time.NewTicker(time.Duration(intervalNum) * time.Second)
	defer ticker.Stop()

	var lastBlockNumber uint64

	for range ticker.C {
		currentBlockNumber := ps.getLatestBlockNumber()

		if currentBlockNumber == 0 || lastBlockNumber == currentBlockNumber {
			log.Print("Cache is all good!")
			continue
		}

		log.Print("Block Number updated, checking cache....")

		cacheKeys := ps.Cache.GetKeys()

		for _, key := range cacheKeys {
			holderInfo, _ := ps.Cache.Get(key)

			if holderInfo.BlockNumber < int(currentBlockNumber) {
				isHolder, err := ps.CheckHolder(holderInfo.HolderAddress)

				if err != nil {
					log.Print("Error checking the holder: ", err.Error())
				}

				log.Print("Updating holder ", key)
				ps.Cache.Set(key, &cache.HolderInfo{
					HolderAddress: key,
					IsHolder:      isHolder,
					BlockNumber:   int(currentBlockNumber),
				})
			}
		}

		lastBlockNumber = currentBlockNumber
	}
}

func (ps *ProxyService) CheckHolder(walletAddress string) (bool, error) {
	return true, nil
}

func (ps *ProxyService) EthSendTransaction(ctx context.Context, args *pb.SendTransactionRequest) (*pb.TransactionReceipt, error) {
	txnHash, err := ps.sendTransaction(args)
	blockNum := ps.getLatestBlockNumber()

	if err != nil {
		return &pb.TransactionReceipt{
			Hash: err.Error(),
		}, nil
	}

	ps.Cache.Set(args.From, &cache.HolderInfo{
		BlockNumber:   int(blockNum),
		HolderAddress: args.From,
		IsHolder:      true,
	})

	return &pb.TransactionReceipt{
		Hash: txnHash,
	}, nil
}
