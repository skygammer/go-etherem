package mymain

import (
	"bytes"
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

	if requestPointer, err :=
		http.NewRequest(
			http.MethodPost,
			`/account`,
			bytes.NewBufferString(`
				{
					"Account":"alice"
				}
			`),
		); err != nil {
		log.Fatal(err)
	} else {
		router.ServeHTTP(responseRecorderPointer, requestPointer)
		assert.Equal(t, http.StatusOK, responseRecorderPointer.Code)
	}

	if requestPointer, err :=
		http.NewRequest(
			http.MethodPost,
			`/account`,
			bytes.NewBufferString(`
			{
				"Account":"bob"
			}
		`),
		); err != nil {
		log.Fatal(err)
	} else {
		router.ServeHTTP(responseRecorderPointer, requestPointer)
		assert.Equal(t, http.StatusOK, responseRecorderPointer.Code)
	}

	if requestPointer, err :=
		http.NewRequest(
			http.MethodPost,
			`/account`,
			bytes.NewBufferString(`
		{
			"Account":"charlie"
		}
	`),
		); err != nil {
		log.Fatal(err)
	} else {
		router.ServeHTTP(responseRecorderPointer, requestPointer)
		assert.Equal(t, http.StatusOK, responseRecorderPointer.Code)
	}

}
