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
	return pkcs7UPad(plaintext)
}

func pkcs7Pad(data []byte) []byte {
	padding := aes.BlockSize - len(data)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// pkcs7UPad removes the PKCS7 padding. A wrong key produces garbage padding,
// so the length is validated instead of trusted: an invalid pad byte means the
// ciphertext was not produced by this key, and slicing on it would panic.
func pkcs7UPad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("ciphertext is empty")
	}
	padLength := int(data[len(data)-1])
	if padLength == 0 || padLength > aes.BlockSize || padLength > len(data) {
		return nil, fmt.Errorf("invalid padding length %d", padLength)
	}
	for _, b := range data[len(data)-padLength:] {
		if int(b) != padLength {
			return nil, fmt.Errorf("invalid padding bytes")
		}
	}
	return data[:len(data)-padLength], nil
}
