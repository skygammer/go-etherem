package mymain

import (
	"errors"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis"
	"github.com/spf13/cobra"
)

var (
	rootCommand = &cobra.Command{
		Args: func(cmd *cobra.Command, args []string) error {

			if len(args) != 1 {
				return errors.New("one argument required")
			} else if len(args[0]) != 8 {
				return errors.New("key for DES encryption/descryption should be in 64 bits")
			} else {
				return nil
			}

		},
		Run: func(cmd *cobra.Command, args []string) {
			desKey = args[0]
		},
	}
)

const (
	AccumulationWalletIndex = iota //歸集錢包
	HotWalletIndex                 //熱錢包
	SystemColdWalletIndex          //系統冷錢包
	BossColdWalletIndex            //boss冷錢包
	WithdrawToWalletIndex          //歸集目的錢包
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
	mnemonic         = `cash village area early major output rubber snow boost amateur option vault`
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
	redisServerUrl           = `localhost:6379`
	userNamespaceConstString = `User`
	userAddressFieldName     = `address`
	userPrivateKeyFieldName  = `private key`
)

// redis var
var (
	desKey string

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
