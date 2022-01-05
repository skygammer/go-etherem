package mymain

import (
	"net/http"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

// API 歸集
func postAccountAccumulationAPI(ginContextPointer *gin.Context) {

	isUndoneChannel <- true

	var (
		parameters struct{}
		httpStatus = http.StatusForbidden
	)

	if !isBindParametersPointerError(ginContextPointer, &parameters) {

		// 将充值地址上的资金归集到指定的地址（第一期使用直接转账的方式，二期需要使用智能合约方式（一次转账可以转出多个地址 ）
		toAddressHex := specialWalletAddressHexes[CollectionIndex]

		keys, _ := redisScan(0, getUserKey(`*`), 0).Val()

		for _, key := range keys {

			if privateKeyBytes, err := redisHGet(key, userPrivateKeyFieldName).Bytes(); err != nil {
				sugaredLogger.Fatal(err)
			} else if privateKeyPointer, err := crypto.HexToECDSA(string(decryptDES(privateKeyBytes, []byte(desKey)))); err != nil {
				sugaredLogger.Fatal(err)
			} else {
				transferBalance(redisHGet(key, userAddressFieldName).Val(), privateKeyPointer, toAddressHex)
			}

		}

		httpStatus = http.StatusOK
	}

	ginContextPointer.Status(httpStatus)
	logAPIRequest(ginContextPointer, parameters, httpStatus)
	<-isUndoneChannel
}
