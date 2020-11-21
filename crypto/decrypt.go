package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

func Decrypt(key []byte, bytes []byte) []byte {
	if len(bytes) == 0 {
		return bytes
	}

	iv := bytes[:16]
	encrypted := bytes[16:]

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encrypted))
	mode.CryptBlocks(decrypted, encrypted)
	return PKCS5Trimming(decrypted)
}

func PKCS5Trimming(b []byte) []byte {
	length := b[len(b)-1]
	return b[:len(b)-int(length)]
}
