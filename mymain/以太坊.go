package mymain

import (
	"crypto/ecdsa"
	"fmt"
	"math"
	"math/big"
	"regexp"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

const (
	AccumulationWalletIndex = iota //歸集錢包
	HotWalletIndex                 //熱錢包
	SystemColdWalletIndex          //系統冷錢包
	BossColdWalletIndex            //boss冷錢包
	WithdrawToWalletIndex          //歸集目的錢包
)

// eth const
const (
	ethHttpServerUrl = `http://localhost:8545`
	ethWsServerUrl   = `ws://localhost:8545`
	accountIndexMax  = 99
	mnemonic         = `holiday poem industry task canoe breeze arm disease coyote kit glimpse avocado`
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

// 判斷地址hex字串是否合法
func isAddressHexStringLegal(addressString string) bool {
	return regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`).MatchString(addressString)
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

// 取得最新兩筆餘額
func getLatestTwoBalances(
	address common.Address,
	blockNumber *big.Int) (lastBalance *big.Int, balance *big.Int, err error) {

	lastBlockNumber := big.NewInt(0).Sub(blockNumber, big.NewInt(1))

	if lastBalance, err =
		ethHttpClientPointer.BalanceAt(
			contextBackground,
			address,
			lastBlockNumber,
		); err != nil {
	} else {

		balance, err =
			ethHttpClientPointer.BalanceAt(
				contextBackground,
				address,
				blockNumber,
			)

	}

	return
}

// 轉帳餘額
func transferBalance(fromAddressHexString string, fromAddressPrivateKeyPointer *ecdsa.PrivateKey, toAddressHexString string) {

	toAddress := common.HexToAddress(toAddressHexString)

	if gasLimit, err := ethHttpClientPointer.EstimateGas(contextBackground, ethereum.CallMsg{To: &toAddress}); err != nil {
		sugaredLogger.Fatal(err)
	} else if amount, err := ethHttpClientPointer.BalanceAt(contextBackground, common.HexToAddress(fromAddressHexString), nil); err != nil {
		sugaredLogger.Fatal(err)
	} else {
		sendTransaction(fromAddressHexString, fromAddressPrivateKeyPointer, toAddressHexString, big.NewInt(0).Add(amount, big.NewInt(int64(-gasLimit))))
	}

}

// 送出交易
func sendTransaction(fromAddressHexString string, fromAddressPrivateKeyPointer *ecdsa.PrivateKey, toAddressHexString string, amount *big.Int) {

	if fromAddressHexString != toAddressHexString && isAddressHexStringLegal(fromAddressHexString) && isAddressHexStringLegal(toAddressHexString) && fromAddressPrivateKeyPointer != nil {

		toAddress := common.HexToAddress(toAddressHexString)

		if gasLimit, err := ethHttpClientPointer.EstimateGas(contextBackground, ethereum.CallMsg{To: &toAddress}); err != nil {
			sugaredLogger.Fatal(err)
		} else if nonce, err := ethHttpClientPointer.PendingNonceAt(contextBackground, common.HexToAddress(fromAddressHexString)); err != nil {
			sugaredLogger.Fatal(err)
		} else if gasPrice, err := ethHttpClientPointer.SuggestGasPrice(contextBackground); err != nil {
			sugaredLogger.Fatal(err)
		} else if chainID, err := ethHttpClientPointer.NetworkID(contextBackground); err != nil {
			sugaredLogger.Fatal(err)
		} else {

			if signedTransactionPointer, err := types.SignTx(types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, nil), types.NewEIP155Signer(chainID), fromAddressPrivateKeyPointer); err != nil {
				sugaredLogger.Fatal(err)
			} else if err := ethHttpClientPointer.SendTransaction(contextBackground, signedTransactionPointer); err != nil {
				sugaredLogger.Fatal(err)
			} else {
				sugaredLogger.Info(fmt.Sprintf(`送出交易 %s`, signedTransactionPointer.Hash().Hex()))
			}

		}

	}

}
