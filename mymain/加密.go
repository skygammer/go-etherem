package mymain

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
)

// 填充
func padding(src []byte, blocksize int) []byte {
	n := len(src)
	padnum := blocksize - n%blocksize
	pad := bytes.Repeat([]byte{byte(padnum)}, padnum)
	dst := append(src, pad...)
	return dst
}

// 反填充
func unpadding(src []byte) []byte {
	n := len(src)
	unpadnum := int(src[n-1])
	dst := src[:n-unpadnum]
	return dst
}

// DES 加密
func encryptDES(src []byte, key []byte) []byte {

	if block, err := des.NewCipher(key); err != nil {
		sugaredLogger.Fatal(err)
	} else {
		src = padding(src, block.BlockSize())
		cipher.NewCBCEncrypter(block, key).CryptBlocks(src, src)
	}

	return src
}

// DES 解密
func decryptDES(src []byte, key []byte) []byte {

	if block, err := des.NewCipher(key); err != nil {
		sugaredLogger.Fatal(err)
	} else {
		cipher.NewCBCDecrypter(block, key).CryptBlocks(src, src)
		src = unpadding(src)
	}

	return src
}
