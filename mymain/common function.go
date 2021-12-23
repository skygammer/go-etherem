package mymain

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

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

// isBindParametersPointerError - 判斷是否綁定參數指標錯誤
func isBindParametersPointerError(ginContextPointer *gin.Context, parametersPointer interface{}) bool {

	var result bool

	if ginContextPointer.Request.Method == http.MethodGet {

		err := ginContextPointer.ShouldBind(parametersPointer)

		result = err != nil

		if result {
			log.Fatal(err)
		}

	} else {

		if rawDataBytes, getRawDataError := ginContextPointer.GetRawData(); getRawDataError != nil {
			log.Fatal(getRawDataError)
		} else {
			ginContextPointer.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawDataBytes))

			err := ginContextPointer.ShouldBindJSON(parametersPointer)

			result = err != nil

			if result {
				log.Fatal(err)
			}

			ginContextPointer.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawDataBytes))

		}

	}

	shouldBindUriError := ginContextPointer.ShouldBindUri(parametersPointer)

	if shouldBindUriError != nil {
		log.Fatal(shouldBindUriError)
	}

	return result || shouldBindUriError != nil
}

// 列印Redis隊列
func printRedisStreams() {

	fmt.Println(`====== redis streams start ======`)

	for _, redisStreamKey := range redisStreamKeys {

		fmt.Println(
			redisClientPointer.XRange(
				redisStreamKey,
				`-`,
				`+`,
			),
		)

	}

	fmt.Println(`====== redis streams end ======`)

}
