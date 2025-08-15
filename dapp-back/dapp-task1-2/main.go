package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	counter "learn-go-ethereum/dapp-task1-2/contract"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

//# 安装 solc（如果本地没有）
//npm install -g solc

// npx solcjs --abi --bin Counter.sol -o build

// go install github.com/ethereum/go-ethereum/cmd/abigen@latest

// abigen --bin=build/Counter_sol_Counter.bin --abi=build/Counter_sol_Counter.abi --pkg=counter --out=Counter.go
func main() {
	rpcURL := "https://sepolia.infura.io/v3/f74b18f9e6e14702a29ea60b9fa5dd89"
	privateKeyHex := "2218488fb96d07282a385c9568445a4ea085e76d5b2ab6b8d83a2428a1321959" // Sepolia 测试网私钥

	//DeployCounterContract(rpcURL, privateKeyHex)
	// 调用合约
	CallCounterContract(rpcURL, privateKeyHex)
}

func CallCounterContract(url string, hex string) {
	// RPC 地址（Infura Sepolia）
	//rpcURL := "https://sepolia.infura.io/v3/<YOUR_INFURA_PROJECT_ID>"

	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}

	// 私钥（Sepolia 测试网）
	privateKey, err := crypto.HexToECDSA(hex)
	if err != nil {
		log.Fatalf("解析私钥失败: %v", err)
	}

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKey)

	// 构造交易签名器
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("获取 nonce 失败: %v", err)
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("获取网络 ID 失败: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("获取 gas price 失败: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("创建交易签名器失败: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // 不转账 ETH
	auth.GasLimit = uint64(300000) // 预估 gas
	auth.GasPrice = gasPrice

	// 合约地址
	contractAddr := common.HexToAddress("0x108Dc28aCE74e5B7a157393f5C48013bCc0e38d0")

	// 实例化合约对象
	counterInstance, err := counter.NewCounter(contractAddr, client)
	if err != nil {
		log.Fatalf("绑定合约失败: %v", err)
	}

	// 调用 increment 方法
	tx, err := counterInstance.Increment(auth)
	if err != nil {
		log.Fatalf("调用 increment 失败: %v", err)
	}
	fmt.Printf("increment 交易已发送: %s\n", tx.Hash().Hex())

	// 读取 count
	count, err := counterInstance.GetCount(&bind.CallOpts{})
	if err != nil {
		log.Fatalf("调用 getCount 失败: %v", err)
	}
	fmt.Printf("当前计数值: %v\n", count)

	//increment 交易已发送: 0xf4946a6a31a1500d8b61a14b7de7d6b7ffbe9f67a61b87814eb0342834da7ab0
	//当前计数值: 2
}

func DeployCounterContract(rpcURL, privateKeyHex string) {
	// 1. 连接到以太坊节点
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("连接节点失败: %v", err)
	}

	// 2. 解析私钥
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("解析私钥失败: %v", err)
	}

	// 3. 获取账户地址
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKey)

	// 4. 获取 nonce、chainID、gasPrice
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("获取 nonce 失败: %v", err)
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("获取网络 ID 失败: %v", err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("获取 gas price 失败: %v", err)
	}

	// 5. 创建交易签名器
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("创建签名器失败: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // 部署不转 ETH
	auth.GasLimit = uint64(300000) // 部署合约的 gas 限制
	auth.GasPrice = gasPrice

	// 6. 部署合约
	address, tx, _, err := counter.DeployCounter(auth, client)
	if err != nil {
		log.Fatalf("部署合约失败: %v", err)
	}

	fmt.Printf("✅ 合约部署成功！\n地址: %s\n交易哈希: %s\n", address.Hex(), tx.Hash().Hex())
	//	✅ 合约部署成功！
	//地址: 0x108Dc28aCE74e5B7a157393f5C48013bCc0e38d0
	//交易哈希: 0xfcddb28f4f09e4c6210685a03fcbec08eb04d2bd0986ade0e2c350655013bd77
}
