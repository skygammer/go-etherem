package mymain

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_postAccountWithdrawalAPI(t *testing.T) {

	initialize()
	redisClientPointer.FlushAll() // 清除之前所有測試資料

	router := setupRouter()
	responseRecorderPointer := httptest.NewRecorder()

	if requestPointer, err :=
		http.NewRequest(
			http.MethodPost,
			`/account/withdrawal/ETH`,
			bytes.NewBufferString(`
				{
					"Address":"0x8d4C1bfc33d20442aA7890196FDf6EFd518eEFE3",
					"Size":1
				}
			`),
		); err != nil {
		log.Fatal(err)
	} else {
		router.ServeHTTP(responseRecorderPointer, requestPointer)
		assert.Equal(t, http.StatusOK, responseRecorderPointer.Code)
	}
}
