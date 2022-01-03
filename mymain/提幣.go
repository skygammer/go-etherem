package mymain

import (
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gin-gonic/gin"
)

// API 提幣
func postAccountWithdrawalAPI(ginContextPointer *gin.Context) {

	isUndoneChannel <- true

	type Parameters struct {
		Address string // 目的地址
		Size    int    // eth值
	}

	var (
		parameters Parameters
		httpStatus = http.StatusForbidden
	)

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
					contextBackground,
					ethereum.CallMsg{
						To: &toAddress,
					},
				); err != nil {
				sugaredLogger.Fatal(err)
			} else if privateKeyPointer, err := fromWalletPointer.PrivateKey(*fromAccountPointer); err != nil {
				sugaredLogger.Fatal(err)
			} else if nonce, err :=
				ethHttpClientPointer.PendingNonceAt(
					contextBackground,
					fromAddress,
				); err != nil {
				sugaredLogger.Fatal(err)
			} else if gasPrice, err :=
				ethHttpClientPointer.SuggestGasPrice(contextBackground); err != nil {
				sugaredLogger.Fatal(err)
			} else if chainID, err := ethHttpClientPointer.NetworkID(contextBackground); err != nil {
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
				} else if err := ethHttpClientPointer.SendTransaction(contextBackground, signedTransactionPointer); err != nil {
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
