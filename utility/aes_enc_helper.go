package utility

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
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

	return fmt.Sprintf("%x", cipherText), nil
}

func AESEncrypt(text []byte, key []byte) (string, *RestError) {
	const ENC_ERROR = "ENC_ERROR"
	c, cErr := aes.NewCipher(key)
	if cErr != nil {
		return "", NewInternalServerError(ENC_ERROR, &cErr)
	}
	gcm, gcmErr := cipher.NewGCM(c)
	if gcmErr != nil {
		return "", NewInternalServerError(ENC_ERROR, &gcmErr)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, rErr := io.ReadFull(rand.Reader, nonce); rErr != nil {
		return "", NewInternalServerError(ENC_ERROR, &rErr)
	}
	result := gcm.Seal(nonce, nonce, text, nil)
	encodedResult := base64.StdEncoding.EncodeToString(result)
	return encodedResult, nil
}
