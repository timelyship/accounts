package utility

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/mergermarket/go-pkcs7"
	"io"
)

func SimpleAESEncrypt(key []byte, unencrypted string) (string, *RestError) {
	plainText := []byte(unencrypted)
	plainText, err := pkcs7.Pad(plainText, aes.BlockSize)
	if err != nil {
		e := fmt.Errorf(`plainText: "%s" has error`, plainText)
		return "", NewUnAuthorizedError("UAE", &e)
	}
	if len(plainText)%aes.BlockSize != 0 {
		err := fmt.Errorf(`plainText: "%s" has the wrong block size`, plainText)
		return "", NewUnAuthorizedError("UAE", &err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", NewUnAuthorizedError("UAE", &err)
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", NewUnAuthorizedError("UAE", &err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainText)

	data := fmt.Sprintf("%x", cipherText)
	return data, nil
}
