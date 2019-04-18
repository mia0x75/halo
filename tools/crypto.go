package tools

import (
	"bytes"
	"crypto/aes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
)

// 对字符串进行sha1 计算
func Sha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return hex.EncodeToString(t.Sum(nil))
}

// 对数据进行md5计算
func MD5(byteMessage []byte) string {
	h := md5.New()
	h.Write(byteMessage)
	return hex.EncodeToString(h.Sum(nil))
}

func HamSha1(data string, key []byte) []byte {
	hmac := hmac.New(sha1.New, key)
	hmac.Write([]byte(data))

	return hmac.Sum(nil)
}

// 使用方式
// plaintext := []byte("C")
// c1 := toolkits.EncryptAES([]byte(plaintext), g.Config().Secret.Crypto)
// p1 := toolkits.DecryptAES(c1, g.Config().Secret.Crypto)
// 如需打印
// fmt.Println(hex.EncodeToString(c1))
// fmt.Println(string(p1))

// 加密
func EncryptAES(plaintext []byte, key string) ([]byte, error) {
	cipher, err := aes.NewCipher([]byte(key[:aes.BlockSize]))
	if err != nil {
		return nil, err
	}

	pt := PKCS7Pad([]byte(plaintext))

	ciphertext := make([]byte, 0)
	text := make([]byte, 16)
	for len(pt) > 0 {
		// 每次运算一个block
		cipher.Encrypt(text, pt)
		pt = pt[aes.BlockSize:]
		ciphertext = append(ciphertext, text...)
	}
	return ciphertext, nil
}

// 解密
func DecryptAES(ciphertext []byte, key string) ([]byte, error) {
	cipher, err := aes.NewCipher([]byte(key[:aes.BlockSize]))
	if err != nil {
		return nil, err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("Need a multiple of the blocksize 16")
	}

	plaintext := make([]byte, 0)
	text := make([]byte, 16)
	for len(ciphertext) > 0 {
		cipher.Decrypt(text, ciphertext)
		ciphertext = ciphertext[aes.BlockSize:]
		plaintext = append(plaintext, text...)
	}
	return PKCS7UPad(plaintext), nil
}

// Padding补全
func PKCS7Pad(data []byte) []byte {
	padding := aes.BlockSize - len(data)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func PKCS7UPad(data []byte) []byte {
	padLength := int(data[len(data)-1])
	return data[:len(data)-padLength]
}
