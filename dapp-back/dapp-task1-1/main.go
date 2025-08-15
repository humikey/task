package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	queryBlock()

	sendTx()
}

func queryBlock() {
	rpcURL := "https://sepolia.infura.io/v3/f74b18f9e6e14702a29ea60b9fa5dd89"

	// 连接到以太坊客户端
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer client.Close()

	// 要查询的区块号
	blockNumber := big.NewInt(4819235) // 你可以改成想查的区块号

	// 查询区块信息
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatalf("查询区块失败: %v", err)
	}

	// 输出区块信息
	fmt.Printf("区块号: %v\n", block.Number().Uint64())
	fmt.Printf("区块哈希: %v\n", block.Hash().Hex())
	fmt.Printf("时间戳: %v (%v)\n", block.Time(), time.Unix(int64(block.Time()), 0))
	fmt.Printf("交易数量: %v\n", len(block.Transactions()))

	//区块号: 4819235
	//区块哈希: 0x5ed3399af8d321287b1a335a73fb7f22b8f2bc4d4eb02b6fb75ec9cc052942c5
	//时间戳: 1701675216 (2023-12-04 15:33:36 +0800 CST)
	//交易数量: 71
}

func sendTx() {
	// Sepolia RPC 地址（替换为你的节点服务，比如 Infura、Alchemy）
	rpcURL := "https://sepolia.infura.io/v3/f74b18f9e6e14702a29ea60b9fa5dd89"

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("连接到节点失败: %v", err)
	}

	// 私钥
	privateKeyHex := "2218488fb96d07282a385c9568445a4ea085e76d5b2ab6b8d83a2428a1321959"

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("解析私钥失败: %v", err)
	}

	// 从私钥获取公钥和地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("公钥类型转换失败")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 获取当前 nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("获取 nonce 失败: %v", err)
	}

	// 交易参数
	toAddress := common.HexToAddress("0x6e9d10214CA7741255b4cA1f19bfCFDB4A46bE71") // 接收方地址
	value := big.NewInt(10000000000000000)                                         // 转账金额: 0.01 ETH (单位: wei)
	gasLimit := uint64(21000)                                                      // 普通转账 gas limit
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("获取 gas price 失败: %v", err)
	}

	// 构造交易
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	// 签名交易
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("获取网络ID失败: %v", err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("签名交易失败: %v", err)
	}

	// 发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("发送交易失败: %v", err)
	}

	fmt.Printf("交易已发送！交易哈希: %s\n", signedTx.Hash().Hex())

	//交易已发送！交易哈希: 0x638d324fc7d7c3d947f6a22214c88d9542d56136cf6299b2e4d6f0e8605a06e5
}
