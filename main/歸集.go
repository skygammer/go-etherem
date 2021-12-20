package main

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gin-gonic/gin"
)

// API 歸集
func postAccountAccumulationAPI(ginContextPointer *gin.Context) {

	type Parameters struct {
		Address string
	}

	var parameters Parameters

	if !isBindParametersPointerError(ginContextPointer, &parameters) {

		// 将充值地址上的资金归集到指定的地址（第一期使用直接转账的方式，二期需要使用智能合约方式（一次转账可以转出多个地址 ）

		privateKeyPointer := accountDatas[AccountWalletIndex].PrivateKeyPointer
		fromAddress := accountDatas[AccountWalletIndex].AccountPointer.Address

		toAddress := common.HexToAddress(parameters.Address)

		gasLimit := uint64(21000) // in units
		data := []byte{}

		if amount, err := ethHttpClientPointer.BalanceAt(context.Background(), fromAddress, nil); err != nil {
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

			transactionPointer :=
				types.NewTransaction(
					nonce,
					toAddress,
					bigIntObject.Add(amount, big.NewInt(int64(-gasLimit))),
					gasLimit,
					gasPrice,
					data,
				)

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
