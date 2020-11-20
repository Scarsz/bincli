package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func Encrypt(key []byte, bytes []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil { panic(err) }

	bytes = PKCS5Padding(bytes)

	encrypted := make([]byte, aes.BlockSize + len(bytes))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encrypted[aes.BlockSize:], bytes)

	return encrypted
}

func EncryptString(key []byte, text string) []byte {
	return Encrypt(key, []byte(text))
}

func PKCS5Padding(b []byte) []byte {
	padding := aes.BlockSize - len(b) % aes.BlockSize
	return append(b, bytes.Repeat([]byte{byte(padding)}, padding)...)
}
