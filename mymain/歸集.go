package mymain

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

// API 歸集
func postAccountAccumulationAPI(ginContextPointer *gin.Context) {

	type Parameters struct {
	}

	var parameters Parameters

	if !isBindParametersPointerError(ginContextPointer, &parameters) {

		// 将充值地址上的资金归集到指定的地址（第一期使用直接转账的方式，二期需要使用智能合约方式（一次转账可以转出多个地址 ）
		toAddress :=
			common.HexToAddress(
				specialWalletAddressHexes[CollectionIndex],
			)

		if gasLimit, err :=
			ethHttpClientPointer.EstimateGas(
				context.Background(),
				ethereum.CallMsg{
					To: &toAddress,
				},
			); err != nil {
			log.Fatal(err)
		} else {

			keys, _ := redisClientPointer.Scan(
				0,
				getUserKey(`*`),
				0,
			).Val()

			for _, key := range keys {

				fromAddressHexString := redisClientPointer.HGet(key, userAddressFieldName).Val()
				fromAddress := common.HexToAddress(fromAddressHexString)

				if privateKeyBytes, err := redisClientPointer.HGet(
					key,
					userPrivateKeyFieldName,
				).Bytes(); err != nil {
					log.Fatal(err)
				} else if privateKeyPointer, err :=
					crypto.HexToECDSA(
						string(
							decryptDES(
								privateKeyBytes,
								[]byte(desKey),
							),
						),
					); err != nil {
					log.Fatal(err)
				} else if amount, err :=
					ethHttpClientPointer.BalanceAt(
						context.Background(),
						fromAddress, nil,
					); err != nil {
					log.Fatal(err)
				} else if nonce, err :=
					ethHttpClientPointer.PendingNonceAt(
						context.Background(),
						fromAddress,
					); err != nil {
					log.Fatal(err)
				} else if gasPrice, err :=
					ethHttpClientPointer.SuggestGasPrice(
						context.Background(),
					); err != nil {
					log.Fatal(err)
				} else if chainID, err :=
					ethHttpClientPointer.NetworkID(
						context.Background(),
					); err != nil {
					log.Fatal(err)
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
						log.Fatal(err)
					} else if err :=
						ethHttpClientPointer.SendTransaction(
							context.Background(),
							signedTransactionPointer,
						); err != nil {
						log.Fatal(err)
					}

				}

			}

		}

	}

}
