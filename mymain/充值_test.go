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

func Test_postAccountDepositAPI(t *testing.T) {

	initialize()

	router := setupRouter()
	responseRecorderPointer := httptest.NewRecorder()

	if walletPointer, accountPointer :=
		getWalletPointerAndAccountPointerByMnemonicStringAndDerivationPathIndex(mnemonic, 0); walletPointer != nil && accountPointer != nil {

		if privateKeyHexString, err := walletPointer.PrivateKeyHex(*accountPointer); err != nil {
			log.Fatal(err)
		} else if requestPointer, err :=
			http.NewRequest(
				http.MethodPost,
				`/account/deposit/ETH`,
				bytes.NewBufferString(
					fmt.Sprintf(`
						{
							"Account":"alice",
							"Address":"%s",
							"PrivateKey":"%s",
							"Size":1
						}
					`,
						accountPointer.Address.Hex(),
						privateKeyHexString,
					),
				),
			); err != nil {
			log.Fatal(err)
		} else {
			router.ServeHTTP(responseRecorderPointer, requestPointer)
			assert.Equal(t, http.StatusOK, responseRecorderPointer.Code)
		}

	}

	if walletPointer, accountPointer :=
		getWalletPointerAndAccountPointerByMnemonicStringAndDerivationPathIndex(mnemonic, 1); walletPointer != nil && accountPointer != nil {

		if privateKeyHexString, err := walletPointer.PrivateKeyHex(*accountPointer); err != nil {
			log.Fatal(err)
		} else if requestPointer, err :=
			http.NewRequest(
				http.MethodPost,
				`/account/deposit/ETH`,
				bytes.NewBufferString(
					fmt.Sprintf(`
					{
						"Account":"alice",
						"Address":"%s",
						"PrivateKey":"%s",
						"Size":1
					}
				`,
						accountPointer.Address.Hex(),
						privateKeyHexString,
					),
				),
			); err != nil {
			log.Fatal(err)
		} else {
			router.ServeHTTP(responseRecorderPointer, requestPointer)
			assert.Equal(t, http.StatusOK, responseRecorderPointer.Code)
		}

	}

	if walletPointer, accountPointer :=
		getWalletPointerAndAccountPointerByMnemonicStringAndDerivationPathIndex(mnemonic, accountIndexMax-4); walletPointer != nil && accountPointer != nil {

		if privateKeyHexString, err := walletPointer.PrivateKeyHex(*accountPointer); err != nil {
			log.Fatal(err)
		} else if requestPointer, err :=
			http.NewRequest(
				http.MethodPost,
				`/account/deposit/ETH`,
				bytes.NewBufferString(
					fmt.Sprintf(`
						{
							"Account":"alice",
							"Address":"%s",
							"PrivateKey":"%s",
							"Size":1
						}
					`,
						accountPointer.Address.Hex(),
						privateKeyHexString,
					),
				),
			); err != nil {
			log.Fatal(err)
		} else {
			router.ServeHTTP(responseRecorderPointer, requestPointer)
			assert.Equal(t, http.StatusOK, responseRecorderPointer.Code)
		}

	}

}
