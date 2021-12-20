package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// isBindParametersPointerError - 判斷是否綁定參數指標錯誤
func isBindParametersPointerError(ginContextPointer *gin.Context, parametersPointer interface{}) bool {

	var result bool

	if ginContextPointer.Request.Method == http.MethodGet {

		err := ginContextPointer.ShouldBind(parametersPointer)

		result = err != nil

		if result {
			log.Fatal(err)
		}

	} else {

		if rawDataBytes, getRawDataError := ginContextPointer.GetRawData(); getRawDataError != nil {
			log.Fatal(getRawDataError)
		} else {
			ginContextPointer.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawDataBytes))

			err := ginContextPointer.ShouldBindJSON(parametersPointer)

			result = err != nil

			if result {
				log.Fatal(err)
			}

			ginContextPointer.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawDataBytes))

		}

	}

	shouldBindUriError := ginContextPointer.ShouldBindUri(parametersPointer)

	if shouldBindUriError != nil {
		log.Fatal(shouldBindUriError)
	}

	return result || shouldBindUriError != nil
}

// 列印Redis隊列
func printRedisList() {

	for _, redisListKey := range redisListKeys {

		fmt.Println(
			redisClientPointer.LRange(
				redisListKey,
				0,
				-1,
			),
		)

	}

}
