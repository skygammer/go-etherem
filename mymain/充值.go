package mymain

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

// API 充值
func postAccountDepositAPI(ginContextPointer *gin.Context) {

	isUndoneChannel <- true

	type Parameters struct {
		Account    string // 帳戶
		Address    string // 來源地址
		PrivateKey string // 來源私鑰
		Size       int    // eth值
	}

	var (
		parameters Parameters
		httpStatus = http.StatusForbidden
	)

	if !isBindParametersPointerError(ginContextPointer, &parameters) {

		if parametersAccount :=
			strings.ToLower(strings.TrimSpace(parameters.Account)); len(parametersAccount) == 0 {
		} else if fromAddressHexString :=
			strings.TrimSpace(parameters.Address); !isAddressHexStringLegal(fromAddressHexString) ||
			isUserAccountAddressHexString(fromAddressHexString) {
		} else if toAddressHexString :=
			redisClientPointer.HGet(
				getUserKey(parametersAccount),
				userAddressFieldName,
			).Val(); fromAddressHexString == toAddressHexString ||
			!isAddressHexStringLegal(toAddressHexString) ||
			!isUserAccountAddressHexString(toAddressHexString) {
		} else if privateKeyPointer, err := crypto.HexToECDSA(parameters.PrivateKey); err != nil {
			sugaredLogger.Fatal(err)
		} else {

			toAddress := common.HexToAddress(toAddressHexString)

			amount :=
				big.NewInt(0).Mul(
					big.NewInt(int64(parameters.Size)),
					weisPerEthBigInt,
				) // in wei (Size eth)

			if gasLimit, err :=
				ethHttpClientPointer.EstimateGas(
					context.Background(),
					ethereum.CallMsg{
						To: &toAddress,
					},
				); err != nil {
				sugaredLogger.Fatal(err)
			} else if nonce, err :=
				ethHttpClientPointer.PendingNonceAt(
					context.Background(),
					common.HexToAddress(fromAddressHexString),
				); err != nil {
				sugaredLogger.Fatal(err)
			} else if gasPrice, err :=
				ethHttpClientPointer.SuggestGasPrice(context.Background()); err != nil {
				sugaredLogger.Fatal(err)
			} else if chainID, err :=
				ethHttpClientPointer.NetworkID(context.Background()); err != nil {
				sugaredLogger.Fatal(err)
			} else {

				transactionPointer :=
					types.NewTransaction(
						nonce,
						toAddress,
						amount,
						gasLimit,
						gasPrice,
						nil,
					)

				if signedTransactionPointer, err :=
					types.SignTx(
						transactionPointer,
						types.NewEIP155Signer(chainID),
						privateKeyPointer,
					); err != nil {
					sugaredLogger.Fatal(err)
				} else if err :=
					ethHttpClientPointer.SendTransaction(
						context.Background(),
						signedTransactionPointer,
					); err != nil {
					sugaredLogger.Fatal(err)
				} else {
					sugaredLogger.Info(
						fmt.Sprintf(
							`送出交易 %s`,
							signedTransactionPointer.Hash().Hex(),
						),
					)
				}

			}

		}

		httpStatus = http.StatusOK

	}

	ginContextPointer.Status(httpStatus)

	logAPIRequest(ginContextPointer, parameters, httpStatus)

	<-isUndoneChannel

}
