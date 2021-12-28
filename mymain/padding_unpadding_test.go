package mymain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_padding_unpadding(t *testing.T) {

	src := []byte(`測試`)

	assert.Equal(
		t,
		src,
		unpadding(padding(src, 5)),
	)

}
