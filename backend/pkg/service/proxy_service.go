package service

import (
	"backend/config"
	"backend/pb"
	"backend/pkg/cache"
	"backend/pkg/utils"
	"context"
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"strconv"
	"time"
)

type ProxyService struct {
	pb.UnimplementedEthProxyServer
	EthClient   *ethclient.Client
	Cache       *cache.LRUCache
	Config      *config.Config
	ProxyNftAbi abi.ABI
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

			if holderInfo.BlockNumber == int(currentBlockNumber) {
				continue
			}

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

		lastBlockNumber = currentBlockNumber
	}
}

func (ps *ProxyService) CheckHolder(walletAddress string) (bool, error) {
	addressType := common.HexToAddress(walletAddress)

	callData, err := ps.ProxyNftAbi.Pack("balanceOf", addressType)

	if err != nil {
		return false, err
	}

	nftAddress := common.HexToAddress(ps.Config.ProxyNftAddress)

	msg := ethereum.CallMsg{
		To:   &nftAddress,
		Data: callData,
	}

	result, err := ps.EthClient.CallContract(context.Background(), msg, nil)

	if err != nil {
		return false, err
	}

	var balance *big.Int
	err = ps.ProxyNftAbi.UnpackIntoInterface(&balance, "balanceOf", result)
	if err != nil {
		log.Print(err.Error())
		return false, err
	}

	return balance.Cmp(big.NewInt(0)) > 0, nil
}

func (ps *ProxyService) AuthorizeHolder(walletAddress string) (bool, error) {
	holderInfo, err := ps.Cache.Get(walletAddress)
	blockNum := ps.getLatestBlockNumber()

	//Doesn't exist in cache or not up to date, check chain
	if err != nil || holderInfo.BlockNumber < int(blockNum) {
		isHolder, err := ps.CheckHolder(walletAddress)

		if err != nil {
			return false, err
		}

		ps.Cache.Set(walletAddress, &cache.HolderInfo{
			BlockNumber:   int(blockNum),
			HolderAddress: walletAddress,
			IsHolder:      isHolder,
		})

		return isHolder, nil
	}

	return holderInfo.IsHolder, nil
}

func (ps *ProxyService) EthSendTransaction(ctx context.Context, args *pb.SendTransactionRequest) (*pb.TransactionReceipt, error) {
	txnHash, err := ps.sendTransaction(args)

	isAuthorized, err := ps.AuthorizeHolder(args.From)

	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	if !isAuthorized {
		return nil, errors.New("Sender doesn't own the access token")
	}

	if err != nil {
		return &pb.TransactionReceipt{
			Hash: err.Error(),
		}, nil
	}

	return &pb.TransactionReceipt{
		Hash: txnHash,
	}, nil
}
