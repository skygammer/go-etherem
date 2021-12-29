package mymain

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isUserAccountAddressHexString(t *testing.T) {

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
	user := `evan`

	if requestPointer, err :=
		http.NewRequest(
			http.MethodPost,
			apiPathString,
			bytes.NewBufferString(
				fmt.Sprintf(
					formatString,
					user,
				),
			),
		); err != nil {
		t.Fatal(err)
	} else {
		router.ServeHTTP(responseRecorderPointer, requestPointer)

		assert.Equal(
			t,
			true,
			isUserAccountAddressHexString(
				redisClientPointer.HGet(
					getUserKey(user),
					userAddressFieldName,
				).Val(),
			),
		)

		assert.Equal(
			t,
			false,
			isUserAccountAddressHexString(
				specialWalletAddressHexes[AccumulationWalletIndex],
			),
		)

	}
}
