package mymain

import (
	"bytes"
	"context"
	"crypto/cipher"
	"crypto/des"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/spf13/cobra"
)

func checkArguments(checks ...cobra.PositionalArgs) cobra.PositionalArgs {

	return func(cmd *cobra.Command, args []string) error {

		for _, check := range checks {

			if err := check(cmd, args); err != nil {
				return err
			}

		}

		if _, err := des.NewCipher([]byte(args[0])); err != nil {
			return err
		} else {
			return nil
		}

	}
}

// 設定路由
func setupRouter() *gin.Engine {

	router := gin.Default()
	router.POST(`/user`, postUserAPI)
	router.POST(`/account/deposit/ETH`, postAccountDepositAPI)
	router.POST(`/account/withdrawal/ETH`, postAccountWithdrawalAPI)
	router.POST(`/accumulation`, postAccountAccumulationAPI)

	return router
}

// 初始化
func initialize() {

	initializeSugaredLogger()

	if err := rootCommand.Execute(); err != nil {
		sugaredLogger.Panic(err)
	} else {

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
			sugaredLogger.Fatal(err)
		} else {
			ethHttpClientPointer = thisEthHttpClientPointer
		}

		// eth websocket 客戶端指標
		if thisEthWebsocketClientPointer, err := ethclient.Dial(ethWsServerUrl); err != nil {
			sugaredLogger.Fatal(err)
		} else {
			ethWebsocketClientPointer = thisEthWebsocketClientPointer
		}

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
		sugaredLogger.Fatal(err)
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
			sugaredLogger.Fatal(err)
		} else if account, err := walletPointer.Derive(derivationPath, false); err != nil {
			sugaredLogger.Fatal(err)
		} else {
			return walletPointer, &account
		}

	}
	return walletPointer, nil

}

// 依據助記詞、推導路徑取得帳戶
func getAccountPointerByMnemonicStringAndDerivationPathString(mnemonicString string, derivationPathString string) *accounts.Account {

	if wallet, err := hdwallet.NewFromMnemonic(mnemonicString); err != nil {
		sugaredLogger.Fatal(err)
	} else if derivationPath, err := hdwallet.ParseDerivationPath(derivationPathString); err != nil {
		sugaredLogger.Fatal(err)
	} else if account, err := wallet.Derive(derivationPath, false); err != nil {
		sugaredLogger.Fatal(err)
	} else {
		return &account
	}

	return nil

}

func isAddressHexStringLegal(addressString string) bool {
	return regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`).MatchString(addressString)
}

// isBindParametersPointerError - 判斷是否綁定參數指標錯誤
func isBindParametersPointerError(ginContextPointer *gin.Context, parametersPointer interface{}) bool {

	var result bool

	if ginContextPointer.Request.Method == http.MethodGet {

		err := ginContextPointer.ShouldBind(parametersPointer)

		result = err != nil

		if result {
			sugaredLogger.Fatal(err)
		}

	} else {

		if rawDataBytes, getRawDataError := ginContextPointer.GetRawData(); getRawDataError != nil {
			sugaredLogger.Fatal(getRawDataError)
		} else {
			ginContextPointer.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawDataBytes))

			err := ginContextPointer.ShouldBindJSON(parametersPointer)

			result = err != nil

			if result {
				sugaredLogger.Fatal(err)
			}

			ginContextPointer.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawDataBytes))

		}

	}

	shouldBindUriError := ginContextPointer.ShouldBindUri(parametersPointer)

	if shouldBindUriError != nil {
		sugaredLogger.Fatal(shouldBindUriError)
	}

	return result || shouldBindUriError != nil
}

// 取得User關鍵字
func getUserKey(userString string) string {

	trimmedUserString := strings.TrimSpace(userString)

	if len(trimmedUserString) > 0 {
		return fmt.Sprintf(
			`%s:%s`,
			userNamespaceConstString,
			trimmedUserString,
		)
	} else {
		return ``
	}

}

// 判斷是否為使用者
func isUser(userString string) bool {

	trimmedUserString := strings.TrimSpace(userString)

	return len(trimmedUserString) > 0 &&
		len(
			redisClientPointer.HKeys(
				getUserKey(trimmedUserString),
			).Val(),
		) > 0
}

// 判斷是否為使用者帳戶hex地址字串
func isUserAccountAddressHexString(addressHexString string) bool {

	if trimmedAddressHexString :=
		strings.TrimSpace(addressHexString); isAddressHexStringLegal(trimmedAddressHexString) {

		keys, _ := redisClientPointer.Scan(
			0,
			getUserKey(`*`),
			0,
		).Val()

		for _, key := range keys {

			if redisClientPointer.HGet(key, userAddressFieldName).Val() ==
				trimmedAddressHexString {
				return true
			}

		}

	}

	return false

}

// 填充
func padding(src []byte, blocksize int) []byte {
	n := len(src)
	padnum := blocksize - n%blocksize
	pad := bytes.Repeat([]byte{byte(padnum)}, padnum)
	dst := append(src, pad...)
	return dst
}

// 反填充
func unpadding(src []byte) []byte {
	n := len(src)
	unpadnum := int(src[n-1])
	dst := src[:n-unpadnum]
	return dst
}

// DES 加密
func encryptDES(src []byte, key []byte) []byte {

	if block, err := des.NewCipher(key); err != nil {
		sugaredLogger.Fatal(err)
	} else {
		src = padding(src, block.BlockSize())
		cipher.NewCBCEncrypter(block, key).CryptBlocks(src, src)
	}

	return src
}

// DES 解密
func decryptDES(src []byte, key []byte) []byte {

	if block, err := des.NewCipher(key); err != nil {
		sugaredLogger.Fatal(err)
	} else {
		cipher.NewCBCDecrypter(block, key).CryptBlocks(src, src)
		src = unpadding(src)
	}

	return src
}

// 取得Transaction關鍵字
func getTransactionKey(hashHexString string) string {

	trimmedHashHexString := strings.TrimSpace(hashHexString)

	if len(trimmedHashHexString) > 0 {
		return fmt.Sprintf(
			`%s:%s`,
			transactionNamespaceConstString,
			trimmedHashHexString,
		)
	} else {
		return ``
	}

}

// 取得最新兩筆餘額
func getLatestTwoBalances(
	address common.Address,
	blockNumber *big.Int) (lastBalance *big.Int, balance *big.Int, err error) {

	lastBlockNumber := big.NewInt(0).Sub(blockNumber, big.NewInt(1))

	if lastBalance, err =
		ethHttpClientPointer.BalanceAt(
			context.Background(),
			address,
			lastBlockNumber,
		); err != nil {
	} else {

		balance, err =
			ethHttpClientPointer.BalanceAt(
				context.Background(),
				address,
				blockNumber,
			)

	}

	return
}
