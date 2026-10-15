package crypto

import (
	"bytes"
	"crypto/aes"
	"fmt"
)

// EncryptAES encrypts plaintext with AES-ECB using PKCS7 padding.
// key must be at least aes.BlockSize bytes; only the first block is used.
func EncryptAES(plaintext []byte, key string) ([]byte, error) {
	cipher, err := aes.NewCipher([]byte(key[:aes.BlockSize]))
	if err != nil {
		return nil, err
	}

	pt := pkcs7Pad([]byte(plaintext))

	ciphertext := make([]byte, 0)
	text := make([]byte, 16)
	for len(pt) > 0 {
		cipher.Encrypt(text, pt)
		pt = pt[aes.BlockSize:]
		ciphertext = append(ciphertext, text...)
	}
	return ciphertext, nil
}

// DecryptAES decrypts AES-ECB ciphertext produced by EncryptAES.
func DecryptAES(ciphertext []byte, key string) ([]byte, error) {
	cipher, err := aes.NewCipher([]byte(key[:aes.BlockSize]))
	if err != nil {
		return nil, err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext must be a multiple of the AES block size")
	}

	plaintext := make([]byte, 0)
	text := make([]byte, 16)
	for len(ciphertext) > 0 {
		cipher.Decrypt(text, ciphertext)
		ciphertext = ciphertext[aes.BlockSize:]
		plaintext = append(plaintext, text...)
	}
	return pkcs7UPad(plaintext), nil
}

func pkcs7Pad(data []byte) []byte {
	padding := aes.BlockSize - len(data)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func pkcs7UPad(data []byte) []byte {
	padLength := int(data[len(data)-1])
	return data[:len(data)-padLength]
}
