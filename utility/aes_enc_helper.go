package utility

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

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
	encodedResult :=base64.StdEncoding.EncodeToString(result)
	return encodedResult, nil
}
