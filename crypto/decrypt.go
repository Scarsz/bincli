package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

func Decrypt(key []byte, bytes []byte) []byte {
	iv := bytes[:16]
	bytes = bytes[16:]

	block, err := aes.NewCipher(key)
	if err != nil { panic(err) }

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(bytes))
	mode.CryptBlocks(decrypted, bytes)

	bytes = PKCS5UnPadding(bytes)

	return bytes
}

func PKCS5UnPadding(b []byte) []byte {
	length := len(b)
	return b[:(length - int(b[length-1]))]
}
