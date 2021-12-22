package mymain

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_postAccountAPI(t *testing.T) {

	initialize()
	redisClientPointer.FlushAll() // 清除之前所有測試資料

	router := setupRouter()
	responseRecorderPointer := httptest.NewRecorder()

	if requestPointer, err :=
		http.NewRequest(
			http.MethodPost,
			`/account`,
			nil,
		); err != nil {
		log.Fatal(err)
	} else {
		router.ServeHTTP(responseRecorderPointer, requestPointer)
		assert.Equal(t, http.StatusOK, responseRecorderPointer.Code)
	}

}
