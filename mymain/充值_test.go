package mymain

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_postAccountDepositAPI(t *testing.T) {

	initialize()

	router := setupRouter()
	responseRecorderPointer := httptest.NewRecorder()
	apiPathString := `/account/deposit/ETH`
	formatString :=
		`
		{
			"Account":"alice",
			"Address":"%s",
			"PrivateKey":"%s",
			"Size":1
		}
	`

	if walletPointer, accountPointer :=
		getWalletPointerAndAccountPointerByMnemonicStringAndDerivationPathIndex(mnemonic, 0); walletPointer != nil && accountPointer != nil {

		if privateKeyHexString, err := walletPointer.PrivateKeyHex(*accountPointer); err != nil {
			t.Fatal(err)
		} else if requestPointer, err :=
			http.NewRequest(
				http.MethodPost,
				apiPathString,
				bytes.NewBufferString(
					fmt.Sprintf(
						formatString,
						accountPointer.Address.Hex(),
						privateKeyHexString,
					),
				),
			); err != nil {
			t.Fatal(err)
		} else {
			router.ServeHTTP(responseRecorderPointer, requestPointer)
			assert.Equal(t, http.StatusOK, responseRecorderPointer.Code)
		}

	}

	if walletPointer, accountPointer :=
		getWalletPointerAndAccountPointerByMnemonicStringAndDerivationPathIndex(mnemonic, 1); walletPointer != nil && accountPointer != nil {

		if privateKeyHexString, err := walletPointer.PrivateKeyHex(*accountPointer); err != nil {
			t.Fatal(err)
		} else if requestPointer, err :=
			http.NewRequest(
				http.MethodPost,
				apiPathString,
				bytes.NewBufferString(
					fmt.Sprintf(
						formatString,
						accountPointer.Address.Hex(),
						privateKeyHexString,
					),
				),
			); err != nil {
			t.Fatal(err)
		} else {
			router.ServeHTTP(responseRecorderPointer, requestPointer)
			assert.Equal(t, http.StatusOK, responseRecorderPointer.Code)
		}

	}

	if walletPointer, accountPointer :=
		getWalletPointerAndAccountPointerByMnemonicStringAndDerivationPathIndex(mnemonic, accountIndexMax-4); walletPointer != nil && accountPointer != nil {

		if privateKeyHexString, err := walletPointer.PrivateKeyHex(*accountPointer); err != nil {
			t.Fatal(err)
		} else if requestPointer, err :=
			http.NewRequest(
				http.MethodPost,
				apiPathString,
				bytes.NewBufferString(
					fmt.Sprintf(
						formatString,
						accountPointer.Address.Hex(),
						privateKeyHexString,
					),
				),
			); err != nil {
			t.Fatal(err)
		} else {
			router.ServeHTTP(responseRecorderPointer, requestPointer)
			assert.Equal(t, http.StatusOK, responseRecorderPointer.Code)
		}

	}

}
