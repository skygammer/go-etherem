package mymain

import (
	"math/big"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// API 提幣
func postAccountWithdrawalAPI(ginContextPointer *gin.Context) {

	isUndoneChannel <- true

	var (
		parameters struct {
			Address string // 目的地址
			Size    int    // eth值
		}
		httpStatus = http.StatusForbidden
	)

	if !isBindParametersPointerError(ginContextPointer, &parameters) {

		if parametersAddress := strings.TrimSpace(parameters.Address); isAddressHexStringLegal(parametersAddress) {

			// 当收到transfer的restful请求时，把热钱包的资金转入到用户指定的钱包地址
			// 离线签名，并广播交易
			fromWalletPointer, fromAccountPointer := getWalletPointerAndAccountPointerByMnemonicStringAndDerivationPathIndex(mnemonic, accountIndexMax-HotWalletIndex)

			if privateKeyPointer, err := fromWalletPointer.PrivateKey(*fromAccountPointer); err != nil {
				sugaredLogger.Fatal(err)
			} else {
				sendTransaction(fromAccountPointer.Address.Hex(), privateKeyPointer, parametersAddress, big.NewInt(0).Mul(big.NewInt(int64(parameters.Size)), weisPerEthBigInt))
			}

		}

		httpStatus = http.StatusOK
	}

	ginContextPointer.Status(httpStatus)
	logAPIRequest(ginContextPointer, parameters, httpStatus)
	<-isUndoneChannel
}
