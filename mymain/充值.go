package mymain

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gin-gonic/gin"
)

// API 充值
func postAccountDepositAPI(ginContextPointer *gin.Context) {

	type Parameters struct {
		Size int
	}

	var parameters Parameters

	if !isBindParametersPointerError(ginContextPointer, &parameters) {

		privateKeyPointer := accountDatas[AccountSourceWalletIndex].PrivateKeyPointer
		fromAddress := accountDatas[AccountSourceWalletIndex].AccountPointer.Address

		toAddress := accountDatas[AccountWalletIndex].AccountPointer.Address

		amount := bigIntObject.Mul(big.NewInt(int64(parameters.Size)), weisPerEthBigInt) // in wei (Size eth)
		gasLimit := uint64(21000)                                                        // in units
		data := []byte{}

		if nonce, err :=
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
