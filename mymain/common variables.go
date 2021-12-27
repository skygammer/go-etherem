package mymain

import (
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis"
)

const (
	AccumulationWalletIndex = iota //歸集錢包
	HotWalletIndex                 //熱錢包
	SystemColdWalletIndex          //系統冷錢包
	BossColdWalletIndex            //boss冷錢包
)

const (
	AccountCreatedIndex = iota // redis對列主題索引
	DepositIndex
	WithdrawIndex
	CollectionIndex
)

var (
	isUndoneChannel = make(chan bool, 100) // channel for is-undone's
)

// eth const
const (
	ethHttpServerUrl = `http://localhost:7545`
	ethWsServerUrl   = `ws://localhost:7545`
	accountIndexMax  = 99
	mnemonic         = `marble spend energy negative egg agree cruise episode trust because hurt mutual`
)

// eth var
var (

	// eth http 客戶端指標
	ethHttpClientPointer *ethclient.Client

	// eth websocket 客戶端指標
	ethWebsocketClientPointer *ethclient.Client

	// 特殊錢包16進位地址
	specialWalletAddressHexes = make([]string, 4)
)

// redis const
const (
	redisServerUrl            = `localhost:6379`
	accountToAddressString    = `account-address`
	addressToPrivateKeyString = `address-privateKey`
)

// redis var
var (

	// redis 客戶端指標
	redisClientPointer = redis.NewClient(
		&redis.Options{
			Addr: redisServerUrl, // redis地址
		},
	)

	// redis列表關鍵字
	redisStreamKeys = []string{
		`account_created`, // 生成用户钱包
		`deposit`,         // 存
		`withdraw`,        // 取
		`collection`,      // 集
	}
)

var (
	bigIntObject     = new(big.Int)
	weisPerEthBigInt = big.NewInt(int64(math.Pow10(18)))
)
