package mymain

import (
	"github.com/ethereum/go-ethereum/ethclient"
)

// 初始化
func initialize() {

	initializeSugaredLogger() // 初始化日志

	if err := rootCommand.Execute(); err != nil {
		sugaredLogger.Panic(err)
	} else {

		for index := range specialWalletAddressHexes {
			specialWalletAddressHexes[index] = getAccountPointerByMnemonicStringAndDerivationPathIndex(mnemonic, accountIndexMax-index).Address.Hex()
		}

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
