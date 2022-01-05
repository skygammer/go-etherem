package mymain

import (
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

// API 充值
func postAccountDepositAPI(ginContextPointer *gin.Context) {

	isUndoneChannel <- true

	var (
		parameters struct {
			Account    string // 帳戶
			Address    string // 來源地址
			PrivateKey string // 來源私鑰
			Size       int    // eth值
		}
		httpStatus = http.StatusForbidden
	)

	if !isBindParametersPointerError(ginContextPointer, &parameters) {

		if parametersAccount := strings.ToLower(strings.TrimSpace(parameters.Account)); len(parametersAccount) == 0 {
		} else if fromAddressHexString := strings.TrimSpace(parameters.Address); isUserAccountAddressHexString(fromAddressHexString) {
		} else if toAddressHexString := redisHGet(getUserKey(parametersAccount), userAddressFieldName).Val(); !isUserAccountAddressHexString(toAddressHexString) {
		} else if privateKeyPointer, err := crypto.HexToECDSA(parameters.PrivateKey); err != nil {
			sugaredLogger.Fatal(err)
		} else {
			sendTransaction(fromAddressHexString, privateKeyPointer, toAddressHexString, big.NewInt(0).Mul(big.NewInt(int64(parameters.Size)), weisPerEthBigInt))
		}

		httpStatus = http.StatusOK
	}

	ginContextPointer.Status(httpStatus)
	logAPIRequest(ginContextPointer, parameters, httpStatus)
	<-isUndoneChannel
}
