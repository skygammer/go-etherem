package mymain

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_postAccountAccumulationAPI(t *testing.T) {

	initialize()

	router := setupRouter()
	responseRecorderPointer := httptest.NewRecorder()
	apiPathString := `/accumulation`
	formatString :=
		`
			{
			}
		`

	if requestPointer, err :=
		http.NewRequest(
			http.MethodPost,
			apiPathString,
			bytes.NewBufferString(
				formatString,
			),
		); err != nil {
		t.Fatal(err)
	} else {
		router.ServeHTTP(responseRecorderPointer, requestPointer)
		assert.Equal(t, http.StatusOK, responseRecorderPointer.Code)
	}
}
