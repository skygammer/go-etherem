package mymain

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
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

	if !isBindParametersPointerError(ginContextPointer, &parameters) {

		if parametersAddress :=
			strings.TrimSpace(
				parameters.Address,
			); isAddressHexStringLegal(parametersAddress) {

			// 当收到transfer的restful请求时，把热钱包的资金转入到用户指定的钱包地址
			// 离线签名，并广播交易

			fromWalletPointer, fromAccountPointer :=
				getWalletPointerAndAccountPointerByMnemonicStringAndDerivationPathIndex(mnemonic, accountIndexMax-HotWalletIndex)
			fromAddress := fromAccountPointer.Address

			toAddress := common.HexToAddress(parameters.Address)

			amount :=
				big.NewInt(0).Mul(big.NewInt(int64(parameters.Size)), weisPerEthBigInt) // in wei (Size eth)

			if gasLimit, err :=
				ethHttpClientPointer.EstimateGas(
					context.Background(),
					ethereum.CallMsg{
						To: &toAddress,
					},
				); err != nil {
				log.Fatal(err)
			} else if privateKeyPointer, err := fromWalletPointer.PrivateKey(*fromAccountPointer); err != nil {
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
					log.Fatal(err)
				} else if err := ethHttpClientPointer.SendTransaction(context.Background(), signedTransactionPointer); err != nil {
					log.Fatal(err)
				} else if receiptPointer, err := ethHttpClientPointer.TransactionReceipt(context.Background(), signedTransactionPointer.Hash()); err != nil {
					log.Fatal(err)
				} else {

					transactionBlockNumber := receiptPointer.BlockNumber
					lastTransactionBlockNumber := big.NewInt(transactionBlockNumber.Int64() - 1)

					if lastFromBalance, err :=
						ethHttpClientPointer.BalanceAt(
							context.Background(),
							fromAddress,
							lastTransactionBlockNumber,
						); err != nil {
						log.Fatal(err)
					} else if fromBalance, err :=
						ethHttpClientPointer.BalanceAt(
							context.Background(),
							fromAddress,
							transactionBlockNumber,
						); err != nil {
						log.Fatal(err)
					} else if lastToBalance, err :=
						ethHttpClientPointer.BalanceAt(
							context.Background(),
							toAddress,
							lastTransactionBlockNumber,
						); err != nil {
						log.Fatal(err)
					} else if toBalance, err :=
						ethHttpClientPointer.BalanceAt(
							context.Background(),
							toAddress,
							transactionBlockNumber,
						); err != nil {
						log.Fatal(err)
					} else {
						log.Println(
							fmt.Sprintf(`
帳戶餘額變化
從 %s : %s -> %s
到 %s : %s -> %s
`,
								fromAddress.Hex(),
								lastFromBalance,
								fromBalance,
								toAddress.Hex(),
								lastToBalance,
								toBalance,
							),
						)
					}

				}

			}

		}

	}

}
