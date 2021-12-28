package mymain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_encryptDES_decryptDES(t *testing.T) {

	src := []byte(`測試`)
	key := []byte(`12345678`)

	assert.Equal(
		t,
		src,
		decryptDES(encryptDES(src, key), key),
	)

}
