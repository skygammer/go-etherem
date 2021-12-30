package mymain

import (
	"context"
	"fmt"
	"math/big"
	"regexp"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

// 判斷地址hex字串是否合法
func isAddressHexStringLegal(addressString string) bool {
	return regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`).MatchString(addressString)
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
