package mymain

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// API 生成用户
func postUserAPI(ginContextPointer *gin.Context) {

	isUndoneChannel <- true

	var (
		parameters struct {
			User string // 新增使用者
		}
		httpStatus = http.StatusForbidden
	)

	if !isBindParametersPointerError(ginContextPointer, &parameters) {

		if parametersUser := strings.ToLower(strings.TrimSpace(parameters.User)); !isUser(parametersUser) {
			const nextUserIndexString = `Next User Index`

			if nextUserIndexNumber, err := redisGet(nextUserIndexString).Int64(); err != nil && err != redis.Nil {
				sugaredLogger.Fatal(err)
			} else if walletPointer, accountPointer := getWalletPointerAndAccountPointerByMnemonicStringAndDerivationPathIndex(mnemonic, int(nextUserIndexNumber)); walletPointer == nil || accountPointer == nil {
			} else if accountPrivateKeyHexString, err := walletPointer.PrivateKeyHex(*accountPointer); err != nil {
				sugaredLogger.Fatal(err)
			} else {
				accountAddressHexString := accountPointer.Address.Hex()
				redisBoolCommandPointer := redisHMSet(getUserKey(parametersUser), map[string]interface{}{userAddressFieldName: accountAddressHexString, userPrivateKeyFieldName: encryptDES([]byte(accountPrivateKeyHexString), []byte(desKey))})
				logRedisBoolCommandPointer(redisBoolCommandPointer)

				if redisBoolCommandPointer.Err() == nil {
					logRedisStatusCommandPointer(redisSet(nextUserIndexString, nextUserIndexNumber+1, 0))
					// 生成account_created消息并发送到队列的account_created主题(redis 中 stream数据)
					logRedisStringCommandPointer(redisXAdd(&redis.XAddArgs{Stream: redisStreamKeys[AccountCreatedIndex], ID: `*`, Values: map[string]interface{}{`user`: parametersUser, `address`: accountAddressHexString}}))
				}
			}
		}
		httpStatus = http.StatusOK
	}
	ginContextPointer.Status(httpStatus)
	logAPIRequest(ginContextPointer, parameters, httpStatus)
	<-isUndoneChannel
}
