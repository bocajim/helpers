package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func encodeBase64(b []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(b))
}

func decodeBase64(b []byte) []byte {
	data, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		panic(err)
	}
	return data
}

func Encrypt(key, text []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Printf(log.Error, "Could not create new cipher using key: " + err.Error())
		return ""
	}
	b := encodeBase64(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Printf(log.Error, "Could not get random string: " + err.Error())
		return ""
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], b)
	return ciphertext
}

func Decrypt(key, text []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Printf(log.Error, "Could not create new cipher using key: " + err.Error())
		return ""
	}
	if len(text) < aes.BlockSize {
		log.Printf(log.Error, "Encrypted text was too short for the blocksize used.")
		return ""
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	return decodeBase64(text)
}
