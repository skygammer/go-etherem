package mymain

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_postAccountWithdrawalAPI(t *testing.T) {

	initialize()

	router := setupRouter()
	responseRecorderPointer := httptest.NewRecorder()

	if requestPointer, err :=
		http.NewRequest(
			http.MethodPost,
			`/account/withdrawal/ETH`,
			bytes.NewBufferString(
				fmt.Sprintf(`
					{
						"Address":"%s",
						"Size":1
					}
					`,
					getAccountPointerByMnemonicStringAndDerivationPathIndex(
						mnemonic,
						accountIndexMax-WithdrawIndex,
					).Address.Hex(),
				),
			),
		); err != nil {
		t.Fatal(err)
	} else {
		router.ServeHTTP(responseRecorderPointer, requestPointer)
		assert.Equal(t, http.StatusOK, responseRecorderPointer.Code)
	}
}
