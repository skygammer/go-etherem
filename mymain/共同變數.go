package mymain

import (
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	AccumulationWalletIndex = iota //歸集錢包
	HotWalletIndex                 //熱錢包
	SystemColdWalletIndex          //系統冷錢包
	BossColdWalletIndex            //boss冷錢包
	WithdrawToWalletIndex          //歸集目的錢包
)

var (
	isUndoneChannel = make(chan bool, 100) // channel for is-undone's
)

// eth const
const (
	ethHttpServerUrl = `http://localhost:7545`
	ethWsServerUrl   = `ws://localhost:7545`
	accountIndexMax  = 99
	mnemonic         = `pulp require estate seed mule snake access elevator afford give bag knife`
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

var (
	weisPerEthBigInt = big.NewInt(int64(math.Pow10(18)))
)
