package mymain

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// API 充值
func getDepositsAPI(ginContextPointer *gin.Context) {

	isUndoneChannel <- true

	type Parameters struct {
		Account   string // 帳戶
		StartTime uint64 // 起始時間
		EndTime   uint64 // 結束時間
	}

	var (
		parameters Parameters
		httpStatus = http.StatusForbidden
	)

	if !isBindParametersPointerError(ginContextPointer, &parameters) {

		parametersAccount :=
			strings.ToLower(strings.TrimSpace(parameters.Account))

		toAddressHexString :=
			redisHGet(
				getUserKey(parametersAccount),
				userAddressFieldName,
			).Val()

		results := make([]map[string]string, 0)

		keys, _ := redisScan(0, getNamespaceKey(redisStreamKeys[DepositIndex], `*`), 0).Val()

		for _, key := range keys {

			if len(parametersAccount) != 0 && toAddressHexString != redisHGet(key, `to`).Val() {
			} else if unixTimestamp, err := redisHGet(key, `time`).Int64(); err != nil {
				sugaredLogger.Fatal(err)
			} else if parameters.StartTime != 0 && parameters.StartTime > uint64(unixTimestamp) {
			} else if parameters.EndTime != 0 && parameters.EndTime < uint64(unixTimestamp) {
			} else {

				results =
					append(results,
						redisClientPointer.HGetAll(contextBackground, key).Val(),
					)
			}

		}

		ginContextPointer.JSON(http.StatusOK, results)

		httpStatus = http.StatusOK

	}

	ginContextPointer.Status(httpStatus)

	logAPIRequest(ginContextPointer, parameters, httpStatus)

	<-isUndoneChannel

}
