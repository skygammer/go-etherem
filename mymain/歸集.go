package mymain

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

// API 歸集
func postAccountAccumulationAPI(ginContextPointer *gin.Context) {

	isUndoneChannel <- true

	type Parameters struct {
	}

	var (
		parameters Parameters
		httpStatus = http.StatusForbidden
	)

	if !isBindParametersPointerError(ginContextPointer, &parameters) {

		// 将充值地址上的资金归集到指定的地址（第一期使用直接转账的方式，二期需要使用智能合约方式（一次转账可以转出多个地址 ）
		toAddress :=
			common.HexToAddress(
				specialWalletAddressHexes[CollectionIndex],
			)

		if gasLimit, err :=
			ethHttpClientPointer.EstimateGas(
				contextBackground,
				ethereum.CallMsg{
					To: &toAddress,
				},
			); err != nil {
			sugaredLogger.Fatal(err)
		} else {

			keys, _ := redisScan(
				0,
				getUserKey(`*`),
				0,
			).Val()

			for _, key := range keys {

				fromAddressHexString :=
					redisHGet(key, userAddressFieldName).Val()
				fromAddress := common.HexToAddress(fromAddressHexString)

				if privateKeyBytes, err := redisHGet(
					key,
					userPrivateKeyFieldName,
				).Bytes(); err != nil {
					sugaredLogger.Fatal(err)
				} else if privateKeyPointer, err :=
					crypto.HexToECDSA(
						string(
							decryptDES(
								privateKeyBytes,
								[]byte(desKey),
							),
						),
					); err != nil {
					sugaredLogger.Fatal(err)
				} else if amount, err :=
					ethHttpClientPointer.BalanceAt(
						contextBackground,
						fromAddress, nil,
					); err != nil {
					sugaredLogger.Fatal(err)
				} else if nonce, err :=
					ethHttpClientPointer.PendingNonceAt(
						contextBackground,
						fromAddress,
					); err != nil {
					sugaredLogger.Fatal(err)
				} else if gasPrice, err :=
					ethHttpClientPointer.SuggestGasPrice(
						contextBackground,
					); err != nil {
					sugaredLogger.Fatal(err)
				} else if chainID, err :=
					ethHttpClientPointer.NetworkID(
						contextBackground,
					); err != nil {
					sugaredLogger.Fatal(err)
				} else {

					transactionPointer :=
						types.NewTransaction(
							nonce,
							toAddress,
							big.NewInt(0).Add(amount, big.NewInt(int64(-gasLimit))),
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
							contextBackground,
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

		}

		httpStatus = http.StatusOK

	}

	ginContextPointer.Status(httpStatus)

	logAPIRequest(ginContextPointer, parameters, httpStatus)

	<-isUndoneChannel

}
