package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"timelyship.com/accounts/config"
)

func main() {
	testEnc()
	config.Init()
	config.Start()
}

func testEnc() {
	fmt.Println("Decryption Program v0.01")

	b64Cypher := "3fmzPYqDv5AWEhRpdy12PvEseQKJjHicxdZBV4LuqsRQX4WAB0wWf1KQM4cCJpnDloOujw9Z6FJxsQqtS2NymtOrB4AjR+xWbCsPoK+aGgjhdj1FobiuagjgtSrBq/giRrCdvMev2em9sSaJ+9k+exym8xM5aPZPVp5q076TWLqQnZgCxM5wszaz7MZEZTcFvx2kB64noIbMXZQ9DwsZ0FdKBDIFATXJUkLd+BK2GZGu08ZtoSRQtrD7nLVn+MIRfEvZgQJIZYX7N+h5c8WWPz98cpddzGYT5Vwc030SKmp605eloQ+YdpGXfrhR1l5bta02sFFyHUJd+ZDz/QyaDw8dbyD3cVqA4yOJ5gdLur024efPLmys+BVk0ezaA+rZc+0rBmgzI5eG+bdHErzMlQdrQaSkijo+AzSCWMFiAkG59Ot82MwmYIT6aw5vpWvmDwNAakvdJuy/5wVvk+ZTfs32Ofxbx7BOa68aa2XHtw88cDxsL+nQYVh8RR2+wxYsCIrnGQ0OOijOQruipvoaFZ2gD33gbEhq89u8xNNxzjMtEeU/1HGOFhdSgvvOLKiKp1snNvWMXE+3dFYiktN1qvduPYy8Yl/lds0PdMytl9LC65DC91/hD/GJQninDwLcGyPL+9kdKrC+lQDNGBg/euzGhhtwGc4b1RvUNuBwkjtdIjAjOmnPwWrc7ml1HKJtwtQs9YBc74GpXZSxz26fo+NFYmGrNk3ND7FUJNpZYgE4xKNzOEUvQaFMA8JddXoseMSEDwbvdNmy/Hf2qwF98eHZV33Bpz9+wBt59Mdf/AVFdsHEaRez+399M71mkbJBPku2/22m4VCb6Jhi1Gf+Tq/E5StAke+MPqMHLL87ZECXH4O3eoQmwxWIxxLNmY/kR13q2Equ5iLA7MsONVLO3NVhJ1rxNfaGMOi9yuZqfG2quhUCe4VMHM5Xwsxr59tseJLVsXYrWsq7jpfGOcb1FDz5oHeXLzeA/V/SzvsoJNH0c1VifCD/GiGIDK/d3O+sB+iBAwJYP/EH28E8O0V5zEYqAxx0UaSGfMobkf8iFSSVhpnckD13nQV4Xtxj1Lyql7279GnWeUN1K3uMaxHXffdNnxDQ1Kl+wbshfSjnQVplitQJBC8xsqMBdYcnm+sydZPHJlvh42YEOYIXcfpog1e9iFRQIMFzo0tTqY6r8WJIQcVuNYUh0TUZ3hHa!"
	ciphertext, _ := base64.StdEncoding.DecodeString(b64Cypher)
	key := []byte("e8f6e90a7edc41a78c04f3ab5f3c5055")

	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		fmt.Println(err)
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(plaintext))
}
