/*
 * 1：为每个用户生成生成用户钱包
 * 2：当监听地址发生新交易时获取通知
 * 3：广播签名交易
 * 4：处理ERC20代币(USDT)的充值
 */

package mymain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
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
	redisListKeys = []string{
		`account_created`, // 生成用户钱包
		`deposit`,         // 存
		`withdraw`,        // 取
		`transfer`,        // 轉
	}
)

var (
	bigIntObject     = new(big.Int)
	weisPerEthBigInt = big.NewInt(int64(math.Pow10(18)))
)

func main() {

	initialize()

	redisClientPointer.FlushAll() // 清除之前所有測試資料

	go func() {
		setupRouter().Run()
	}()

	// 1：为每个用户生成生成用户钱包
	// 2：当监听地址发生新交易时获取通知
	// 3：广播签名交易
	// 4：处理ERC20代币(USDT)的充值

	headerChannel := make(chan *types.Header)

	if subscription, err :=
		ethWebsocketClientPointer.SubscribeNewHead(
			context.Background(),
			headerChannel,
		); err != nil {
		log.Fatal(err)
	} else {

		signalChannel := make(chan os.Signal, 1)                    // channel for interrupt
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM) // notify interrupts
		isDoneChannel := make(chan bool, 100)

		for {

			select {

			case err := <-subscription.Err():

				if err != nil {
					log.Fatal(err)
				}

			case header := <-headerChannel:

				// 监听新区块，获取区块中的全部交易
				// 过滤掉与钱包地址无关的交易
				// 将每个相关的交易都发往队列

				if block, err := ethWebsocketClientPointer.BlockByHash(context.Background(), header.Hash()); err != nil {
					log.Fatal(err)
				} else {

					for _, tx := range block.Transactions() {

						toAddress := tx.To()

						if fromAddress, err :=
							types.Sender(types.NewEIP2930Signer(tx.ChainId()), tx); err != nil {
							log.Fatal(err)
						} else {

							valueInETHString := bigIntObject.Div(tx.Value(), weisPerEthBigInt).String()

							fromAddressHex := fromAddress.Hex()

							toAddressHex := toAddress.Hex()

							if fromAddress ==
								accountDatas[AccountSourceWalletIndex].AccountPointer.Address &&
								*toAddress == accountDatas[AccountWalletIndex].AccountPointer.Address {

								// 生成transfer消息并发送到队列的deposit主题(redis 中 stream数据)
								redisClientPointer.RPush(
									redisListKeys[DepositIndex],
									fmt.Sprintf(
										`帳戶 %v 充值 %v ETH`,
										toAddressHex,
										valueInETHString,
									),
								)

							} else if fromAddress ==
								accountDatas[HotWalletIndex].AccountPointer.Address {

								// 生成transfer消息并发送到队列的withdraw主题(redis 中 stream数据)
								redisClientPointer.RPush(
									redisListKeys[WithdrawIndex],
									fmt.Sprintf(
										`帳戶 %v 提幣 %v ETH 到 帳戶 %v`,
										fromAddressHex,
										valueInETHString,
										toAddressHex,
									),
								)

							} else if fromAddress ==
								accountDatas[AccountWalletIndex].AccountPointer.Address {

								// 生成transfer消息并发送到队列的collection主题(redis 中 stream数据)
								redisClientPointer.RPush(
									redisListKeys[CollectionIndex],
									fmt.Sprintf(
										`帳戶 %v 歸集 %v ETH 到 帳戶 %v`,
										fromAddressHex,
										valueInETHString,
										toAddressHex,
									),
								)

							}

							printRedisList()

						}

					}

				}

				isDoneChannel <- true

			case signal := <-signalChannel:
				log.Println(signal)
				close(headerChannel)
				subscription.Unsubscribe()

				for len(isDoneChannel) != 0 {
					<-isDoneChannel
				}

				return
			}

		}

	}

}

// 設定路由
func setupRouter() *gin.Engine {

	router := gin.Default()
	router.POST(`/account`, postAccountAPI)
	router.POST(`/account/deposit/ETH`, postAccountDepositAPI)
	router.POST(`/account/withdrawal/ETH`, postAccountWithdrawalAPI)
	router.POST(`/account/accumulation`, postAccountAccumulationAPI)

	return router
}

// 初始化
func initialize() {

	// 用戶來源錢包
	accountDatas[AccountSourceWalletIndex] =
		AccountData{
			// Mnemonic:         `resist apart surround clinic ivory arrow decide company nut gentle broom taste`,
			DerivationPathIndex: AccountSourceWalletIndex,
			PrivateKeyString:    `c0effed41fc8f8539e0fc88fadc27a13b8f669dc21ebcea6a0d75341e811f0c1`,
		}

	// 用戶錢包
	accountDatas[AccountWalletIndex] =
		AccountData{
			// Mnemonic:         `resist apart surround clinic ivory arrow decide company nut gentle broom taste`,
			DerivationPathIndex: AccountWalletIndex,
			PrivateKeyString:    `aeabbb9d31fe1be71d93d62f9cb1e6e4bece078ca408ed3dd6af31cea18b46c9`,
		}

	// 歸集錢包
	accountDatas[AccumulationWalletIndex] =
		AccountData{
			// Mnemonic:         `essence impulse number there double bottom phrase wink foster stamp affair three`,
			DerivationPathIndex: AccountSourceWalletIndex,
			PrivateKeyString:    `d7b173f29b262428fd8d732014aa6d81bc4cc52059d72be1142c73080d63b3e8`,
		}

	//熱錢包
	accountDatas[HotWalletIndex] =
		AccountData{
			// Mnemonic:         `laundry antenna erupt galaxy tattoo test okay laptop endless reject trend planet`,
			DerivationPathIndex: HotWalletIndex,
			PrivateKeyString:    `4119ac2270a283f5df917fde1b79b7602f0e02406bf9db0ceeb4940bb62e7423`,
		}

	// 系統冷錢包
	accountDatas[SystemColdWalletIndex] =
		AccountData{
			// Mnemonic:         `scheme inject column require story gown rabbit escape movie forward hybrid place`,
			DerivationPathIndex: SystemColdWalletIndex,
			PrivateKeyString:    `b4cf32bc28145bf754720276d03752a46f85f72d2436ba7375209afca258d3e3`,
		}

	// boss冷錢包
	accountDatas[BossColdWalletIndex] =
		AccountData{
			// Mnemonic:         `grace bring evoke proud endless figure convince ready acid afford plastic disagree`,
			DerivationPathIndex: BossColdWalletIndex,
			PrivateKeyString:    `8b06558092942ad8e1a90b43360babbb1cd8416e0348f53458464a871877528d`,
		}

	for index, value := range accountDatas {

		// 助記詞
		value.Mnemonic = `length frame sorry say hockey simple tired document sing mail melt estate`

		// 私鑰指標
		value.PrivateKeyPointer = getPrivateKeyPointerFromPrivateKeyString(value.PrivateKeyString)

		// 公鑰指標
		value.PublicKeyPointer = getPublicKeyPointerFromPrivateKeyString(value.PrivateKeyString)

		// 帳戶指標
		value.AccountPointer =
			getAccountPointerByMnemonicStringAndDerivationPathIndex(
				value.Mnemonic,
				value.DerivationPathIndex,
			)

		accountDatas[index] = value
	}

	// eth http 客戶端指標
	if thisEthHttpClientPointer, err := ethclient.Dial(`http://localhost:7545`); err != nil {
		log.Fatal(err)
	} else {
		ethHttpClientPointer = thisEthHttpClientPointer
	}

	// eth websocket 客戶端指標
	if thisEthWebsocketClientPointer, err := ethclient.Dial(`ws://localhost:7545`); err != nil {
		log.Fatal(err)
	} else {
		ethWebsocketClientPointer = thisEthWebsocketClientPointer
	}

}

// 依據助記詞取得預設帳戶
func getDefaultAccountPointerByMnemonicStringAndDerivationPathString(mnemonicString string) *accounts.Account {

	return getAccountPointerByMnemonicStringAndDerivationPathIndex(
		mnemonicString,
		0,
	)

}

// 依據助記詞、推導路徑索引取得帳戶
func getAccountPointerByMnemonicStringAndDerivationPathIndex(mnemonicString string, derivationPathIndex int) *accounts.Account {

	return getAccountPointerByMnemonicStringAndDerivationPathString(
		mnemonicString,
		fmt.Sprintf(
			`m/44'/60'/0'/0/%d`,
			derivationPathIndex,
		),
	)

}

// 依據助記詞、推導路徑取得帳戶
func getAccountPointerByMnemonicStringAndDerivationPathString(mnemonicString string, derivationPathString string) *accounts.Account {

	if wallet, err := hdwallet.NewFromMnemonic(mnemonicString); err != nil {
		log.Fatal(err)
	} else if derivationPath, err := hdwallet.ParseDerivationPath(derivationPathString); err != nil {
		log.Fatal(err)
	} else if account, err := wallet.Derive(derivationPath, false); err != nil {
		log.Fatal(err)
	} else {
		return &account
	}

	return nil

}

func getPrivateKeyPointerFromPrivateKeyString(privateKeyString string) *ecdsa.PrivateKey {

	if privateKeyPointer, err := crypto.HexToECDSA(privateKeyString); err != nil {
		log.Fatal(err)
		return nil
	} else {
		return privateKeyPointer
	}

}

func getPublicKeyPointerFromPrivateKeyString(privateKeyString string) *ecdsa.PublicKey {

	if privateKeyPointer := getPrivateKeyPointerFromPrivateKeyString(privateKeyString); privateKeyPointer == nil {
		return nil
	} else if publicKeyPointer, ok := privateKeyPointer.Public().(*ecdsa.PublicKey); !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		return nil
	} else {
		return publicKeyPointer
	}

}
