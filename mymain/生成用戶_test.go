package mymain

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_postAccountAPI(t *testing.T) {

	initialize()

	router := setupRouter()
	responseRecorderPointer := httptest.NewRecorder()
	apiPathString := `/user`
	formatString :=
		`
			{
				"User":"%s"
			}
		`

	if requestPointer, err :=
		http.NewRequest(
			http.MethodPost,
			apiPathString,
			bytes.NewBufferString(
				fmt.Sprintf(
					formatString,
					`alice`,
				),
			),
		); err != nil {
		log.Fatal(err)
	} else {
		router.ServeHTTP(responseRecorderPointer, requestPointer)
		assert.Equal(t, http.StatusOK, responseRecorderPointer.Code)
	}

	if requestPointer, err :=
		http.NewRequest(
			http.MethodPost,
			apiPathString,
			bytes.NewBufferString(
				fmt.Sprintf(
					formatString,
					`bob`,
				),
			),
		); err != nil {
		log.Fatal(err)
	} else {
		router.ServeHTTP(responseRecorderPointer, requestPointer)
		assert.Equal(t, http.StatusOK, responseRecorderPointer.Code)
	}

	if requestPointer, err :=
		http.NewRequest(
			http.MethodPost,
			apiPathString,
			bytes.NewBufferString(
				fmt.Sprintf(
					formatString,
					`charlie`,
				),
			),
		); err != nil {
		log.Fatal(err)
	} else {
		router.ServeHTTP(responseRecorderPointer, requestPointer)
		assert.Equal(t, http.StatusOK, responseRecorderPointer.Code)
	}

}
