package mymain

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isUser(t *testing.T) {

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
	user := `david`

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
			isUser(
				user,
			),
		)

		assert.Equal(
			t,
			false,
			isUser(
				`#######`,
			),
		)

	}
}
