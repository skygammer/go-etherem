package mymain

import (
	"context"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

// API 充值
func postAccountDepositAPI(ginContextPointer *gin.Context) {

	type Parameters struct {
		Account    string
		Address    string
		PrivateKey string
		Size       int
	}

	var parameters Parameters

	if !isBindParametersPointerError(ginContextPointer, &parameters) {

		if parametersAccount :=
			strings.ToLower(strings.TrimSpace(parameters.Account)); len(parametersAccount) == 0 {
		} else if fromAddressHexString :=
			strings.TrimSpace(parameters.Address); !isAddressHexStringLegal(fromAddressHexString) ||
			len(
				redisClientPointer.HGet(
					addressToPrivateKeyString,
					fromAddressHexString,
				).Val(),
			) > 0 {
		} else if toAddressHexString :=
			redisClientPointer.HGet(
				accountToAddressString, parametersAccount,
			).Val(); fromAddressHexString == toAddressHexString ||
			!isAddressHexStringLegal(toAddressHexString) {
		} else if privateKeyPointer, err := crypto.HexToECDSA(parameters.PrivateKey); err != nil {
			log.Fatal(err)
		} else {

			amount :=
				bigIntObject.Mul(big.NewInt(int64(parameters.Size)), weisPerEthBigInt) // in wei (Size eth)
			gasLimit := uint64(21000) // in units
			data := []byte{}

			if nonce, err :=
				ethHttpClientPointer.PendingNonceAt(
					context.Background(),
					common.HexToAddress(fromAddressHexString),
				); err != nil {
				log.Fatal(err)
			} else if gasPrice, err :=
				ethHttpClientPointer.SuggestGasPrice(context.Background()); err != nil {
				log.Fatal(err)
			} else if chainID, err :=
				ethHttpClientPointer.NetworkID(context.Background()); err != nil {
				log.Fatal(err)
			} else {

				transactionPointer := types.NewTransaction(nonce, common.HexToAddress(toAddressHexString), amount, gasLimit, gasPrice, data)

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

}
