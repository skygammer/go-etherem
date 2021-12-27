package mymain

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

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

	// 歸集錢包
	specialWalletAddressHexes[AccumulationWalletIndex] =
		getAccountPointerByMnemonicStringAndDerivationPathIndex(
			mnemonic,
			accountIndexMax-AccountCreatedIndex,
		).Address.Hex()

	//熱錢包
	specialWalletAddressHexes[HotWalletIndex] =
		getAccountPointerByMnemonicStringAndDerivationPathIndex(
			mnemonic, accountIndexMax-HotWalletIndex,
		).Address.Hex()

	// 系統冷錢包
	specialWalletAddressHexes[SystemColdWalletIndex] =
		getAccountPointerByMnemonicStringAndDerivationPathIndex(
			mnemonic,
			accountIndexMax-SystemColdWalletIndex,
		).Address.Hex()

	// boss冷錢包
	specialWalletAddressHexes[BossColdWalletIndex] =
		getAccountPointerByMnemonicStringAndDerivationPathIndex(
			mnemonic,
			accountIndexMax-BossColdWalletIndex,
		).Address.Hex()

	// eth http 客戶端指標
	if thisEthHttpClientPointer, err := ethclient.Dial(ethHttpServerUrl); err != nil {
		log.Fatal(err)
	} else {
		ethHttpClientPointer = thisEthHttpClientPointer
	}

	// eth websocket 客戶端指標
	if thisEthWebsocketClientPointer, err := ethclient.Dial(ethWsServerUrl); err != nil {
		log.Fatal(err)
	} else {
		ethWebsocketClientPointer = thisEthWebsocketClientPointer
	}

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

// 依據助記詞取得錢包指標
func getWalletPointerByMnemonicString(mnemonicString string) *hdwallet.Wallet {

	walletPointer, err := hdwallet.NewFromMnemonic(mnemonicString)

	if err != nil {
		log.Fatal(err)
	}

	return walletPointer
}

// 依據助記詞、推導路徑索引取得錢包指標、帳戶指標
func getWalletPointerAndAccountPointerByMnemonicStringAndDerivationPathIndex(mnemonicString string, derivationPathIndex int) (*hdwallet.Wallet, *accounts.Account) {

	return getWalletPointerAndAccountPointerByMnemonicStringAndDerivationPathString(
		mnemonicString,
		fmt.Sprintf(
			`m/44'/60'/0'/0/%d`,
			derivationPathIndex,
		),
	)

}

// 依據助記詞、推導路徑取得錢包指標、帳戶指標
func getWalletPointerAndAccountPointerByMnemonicStringAndDerivationPathString(mnemonicString string, derivationPathString string) (*hdwallet.Wallet, *accounts.Account) {

	walletPointer := getWalletPointerByMnemonicString(mnemonicString)

	if walletPointer != nil {

		if derivationPath, err := hdwallet.ParseDerivationPath(derivationPathString); err != nil {
			log.Fatal(err)
		} else if account, err := walletPointer.Derive(derivationPath, false); err != nil {
			log.Fatal(err)
		} else {
			return walletPointer, &account
		}

	}
	return walletPointer, nil

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

func isAddressHexStringLegal(addressString string) bool {
	return regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`).MatchString(addressString)
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

	for _, redisStreamKey := range redisStreamKeys {

		log.Println(
			redisClientPointer.XRange(
				redisStreamKey,
				`-`,
				`+`,
			),
		)

	}

}
