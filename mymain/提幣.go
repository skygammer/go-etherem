package mymain

import (
	"context"
	"log"
	"math/big"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gin-gonic/gin"
)

// API 提幣
func postAccountWithdrawalAPI(ginContextPointer *gin.Context) {

	type Parameters struct {
		Address string
		Size    int
	}

	var parameters Parameters

	if !isBindParametersPointerError(ginContextPointer, &parameters) &&
		regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`).MatchString(parameters.Address) {

		// 当收到transfer的restful请求时，把热钱包的资金转入到用户指定的钱包地址
		// 离线签名，并广播交易

		privateKeyPointer := accountDatas[HotWalletIndex].PrivateKeyPointer
		fromAddress := accountDatas[HotWalletIndex].AccountPointer.Address

		toAddress := common.HexToAddress(parameters.Address)

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
