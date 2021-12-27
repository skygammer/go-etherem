package mymain

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gin-gonic/gin"
)

// API 提幣
func postAccountWithdrawalAPI(ginContextPointer *gin.Context) {

	type Parameters struct {
		Address string // 目的地址
		Size    int    // eth值
	}

	var parameters Parameters

	if !isBindParametersPointerError(ginContextPointer, &parameters) &&
		isAddressHexStringLegal(parameters.Address) {

		// 当收到transfer的restful请求时，把热钱包的资金转入到用户指定的钱包地址
		// 离线签名，并广播交易

		fromWalletPointer, fromAccountPointer := getWalletPointerAndAccountPointerByMnemonicStringAndDerivationPathIndex(mnemonic, accountIndexMax-HotWalletIndex)
		fromAddress := fromAccountPointer.Address

		toAddress := common.HexToAddress(parameters.Address)

		amount := bigIntObject.Mul(big.NewInt(int64(parameters.Size)), weisPerEthBigInt) // in wei (Size eth)
		gasLimit := uint64(21000)                                                        // in units
		data := []byte{}

		if privateKeyPointer, err := fromWalletPointer.PrivateKey(*fromAccountPointer); err != nil {
			log.Fatal(err)
		} else if nonce, err :=
			ethHttpClientPointer.PendingNonceAt(
				context.Background(),
				fromAddress,
			); err != nil {
			log.Fatal(err)
		} else if gasPrice, err :=
			ethHttpClientPointer.SuggestGasPrice(context.Background()); err != nil {
			log.Fatal(err)
		} else if chainID, err := ethHttpClientPointer.NetworkID(context.Background()); err != nil {
			log.Fatal(err)
		} else {

			transactionPointer := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, data)

			if signedTransactionPointer, err :=
				types.SignTx(
					transactionPointer,
					types.NewEIP155Signer(chainID),
					privateKeyPointer,
				); err != nil {
				log.Fatal(err)
			} else if err := ethHttpClientPointer.SendTransaction(context.Background(), signedTransactionPointer); err != nil {
				log.Fatal(err)
			}

		}

	}

}
