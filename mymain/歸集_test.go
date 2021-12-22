package mymain

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_postAccountAccumulationAPI(t *testing.T) {
	initialize()
	router := setupRouter()
	responseRecorderPointer := httptest.NewRecorder()

	if requestPointer, err :=
		http.NewRequest(
			http.MethodPost,
			`/account/withdrawal/ETH`,
			bytes.NewBufferString(`
				{
					"Address":"0x7b2055Bb6c42704980Fb48064Bb3E24C292b9ED0"
				}
			`),
		); err != nil {
		log.Fatal(err)
	} else {
		router.ServeHTTP(responseRecorderPointer, requestPointer)
		assert.Equal(t, http.StatusOK, responseRecorderPointer.Code)
	}
}
