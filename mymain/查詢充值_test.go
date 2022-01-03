package mymain

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_getDepositsAPI(t *testing.T) {

	initialize()

	router := setupRouter()
	responseRecorderPointer := httptest.NewRecorder()
	apiPathString := `/deposits`

	params := make(url.Values)

	if requestPointer, err :=
		http.NewRequest(
			http.MethodGet,
			strings.Join([]string{apiPathString, params.Encode()}, `?`),
			nil,
		); err != nil {
		t.Fatal(err)
	} else {
		router.ServeHTTP(responseRecorderPointer, requestPointer)
		result, err := ioutil.ReadAll(responseRecorderPointer.Body)
		assert.Equal(t, nil, err)
		t.Log(string(result))
		assert.Equal(t, http.StatusOK, responseRecorderPointer.Code)
	}

	params.Add(`Account`, `alice`)

	if requestPointer, err :=
		http.NewRequest(
			http.MethodGet,
			strings.Join([]string{apiPathString, params.Encode()}, `?`),
			nil,
		); err != nil {
		t.Fatal(err)
	} else {
		router.ServeHTTP(responseRecorderPointer, requestPointer)
		result, err := ioutil.ReadAll(responseRecorderPointer.Body)
		assert.Equal(t, nil, err)
		t.Log(string(result))
		assert.Equal(t, http.StatusOK, responseRecorderPointer.Code)
	}

	params.Del(`Account`)
	params.Add(`StartTime`, fmt.Sprint(time.Now().Unix()))

	if requestPointer, err :=
		http.NewRequest(
			http.MethodGet,
			strings.Join([]string{apiPathString, params.Encode()}, `?`),
			nil,
		); err != nil {
		t.Fatal(err)
	} else {
		router.ServeHTTP(responseRecorderPointer, requestPointer)
		result, err := ioutil.ReadAll(responseRecorderPointer.Body)
		assert.Equal(t, nil, err)
		t.Log(string(result))
		assert.Equal(t, http.StatusOK, responseRecorderPointer.Code)
	}

	params.Add(`Account`, `alice`)
	params.Del(`StartTime`)
	params.Add(`EndTime`, fmt.Sprint(time.Now().Unix()))

	if requestPointer, err :=
		http.NewRequest(
			http.MethodGet,
			strings.Join([]string{apiPathString, params.Encode()}, `?`),
			nil,
		); err != nil {
		t.Fatal(err)
	} else {
		router.ServeHTTP(responseRecorderPointer, requestPointer)
		result, err := ioutil.ReadAll(responseRecorderPointer.Body)
		assert.Equal(t, nil, err)
		t.Log(string(result))
		assert.Equal(t, http.StatusOK, responseRecorderPointer.Code)
	}

	params.Add(`StartTime`, fmt.Sprint(time.Now().Unix()))
	params.Add(`EndTime`, fmt.Sprint(time.Now().Unix()))

	if requestPointer, err :=
		http.NewRequest(
			http.MethodGet,
			strings.Join([]string{apiPathString, params.Encode()}, `?`),
			nil,
		); err != nil {
		t.Fatal(err)
	} else {
		router.ServeHTTP(responseRecorderPointer, requestPointer)
		result, err := ioutil.ReadAll(responseRecorderPointer.Body)
		assert.Equal(t, nil, err)
		t.Log(string(result))
		assert.Equal(t, http.StatusOK, responseRecorderPointer.Code)
	}

}
