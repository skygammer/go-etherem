package mymain

import (
	"crypto/ecdsa"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis"
)

// 用戶資料
type AccountData struct {
	Mnemonic            string            // 助記詞
	DerivationPathIndex int               // 導出路徑編號
	PrivateKeyString    string            // 私鑰字串
	PrivateKeyPointer   *ecdsa.PrivateKey // 私鑰指標
	PublicKeyPointer    *ecdsa.PublicKey  // 公鑰指標
	AccountPointer      *accounts.Account // 帳戶指標
}

const (
	AccountSourceWalletIndex = iota // 用戶來源錢包
	AccountWalletIndex              //用戶錢包
	AccumulationWalletIndex         //歸集錢包
	HotWalletIndex                  //熱錢包
	SystemColdWalletIndex           //系統冷錢包
	BossColdWalletIndex             //boss冷錢包
)

const (
	AccountCreatedIndex = iota // redis對列主題索引
	DepositIndex
	WithdrawIndex
	CollectionIndex
)

var (

	// eth http 客戶端指標
	ethHttpClientPointer *ethclient.Client

	// eth websocket 客戶端指標
	ethWebsocketClientPointer *ethclient.Client

	// redis 客戶端指標
	redisClientPointer = redis.NewClient(
		&redis.Options{
			Addr:     `localhost:6379`, // redis地址
			Password: ``,               // redis密码，没有则留空
			DB:       0,                // 默认数据库，默认是0
		},
	)

	// 編號-帳戶資料對應
	accountDatas = make([]AccountData, 6)

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
