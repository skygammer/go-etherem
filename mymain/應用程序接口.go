package mymain

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 設定路由
func setupRouter() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST(`/user`, postUserAPI)
	router.POST(`/account/deposit/ETH`, postAccountDepositAPI)
	router.POST(`/account/withdrawal/ETH`, postAccountWithdrawalAPI)
	router.POST(`/accumulation`, postAccountAccumulationAPI)
	router.GET(`/deposits`, getDepositsAPI)

	return router
}

// isBindParametersPointerError - 判斷是否綁定參數指標錯誤
func isBindParametersPointerError(ginContextPointer *gin.Context, parametersPointer interface{}) bool {

	var result bool

	if ginContextPointer.Request.Method == http.MethodGet {

		err := ginContextPointer.ShouldBind(parametersPointer)

		result = err != nil

		if result {
			sugaredLogger.Fatal(err)
		}

	} else {

		if rawDataBytes, getRawDataError := ginContextPointer.GetRawData(); getRawDataError != nil {
			sugaredLogger.Fatal(getRawDataError)
		} else {
			ginContextPointer.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawDataBytes))

			err := ginContextPointer.ShouldBindJSON(parametersPointer)

			result = err != nil

			if result {
				sugaredLogger.Fatal(err)
			}

			ginContextPointer.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawDataBytes))

		}

	}

	shouldBindUriError := ginContextPointer.ShouldBindUri(parametersPointer)

	if shouldBindUriError != nil {
		sugaredLogger.Fatal(shouldBindUriError)
	}

	return result || shouldBindUriError != nil
}

// 紀錄API請求
func logAPIRequest(ginContextPointer *gin.Context, parameters interface{}, httpStatus int) {

	if ginContextPointer != nil && ginContextPointer.Request != nil {

		request := ginContextPointer.Request

		sugaredLogger.Info(
			fmt.Sprintf(
				`%s %s %s %+v %d %s`,
				request.RemoteAddr,
				request.Method,
				request.URL,
				parameters,
				httpStatus,
				http.StatusText(httpStatus),
			),
		)

	}

}
