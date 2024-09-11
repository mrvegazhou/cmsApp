package AES

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	b64 "encoding/base64"
	"fmt"
	"io"
)

// Encrypts text with the passphrase
func Encrypt(text string, passphrase string) string {
	salt := make([]byte, 8)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		panic(err.Error())
	}

	key, iv := __DeriveKeyAndIv(passphrase, string(salt))

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	pad := __PKCS7Padding([]byte(text), block.BlockSize())
	ecb := cipher.NewCBCEncrypter(block, []byte(iv))
	encrypted := make([]byte, len(pad))
	ecb.CryptBlocks(encrypted, pad)

	return b64.StdEncoding.EncodeToString([]byte("S_" + string(salt) + string(encrypted)))
}

// Decrypts encrypted text with the passphrase
func Decrypt(encrypted string, passphrase string) string {
	ct, _ := b64.StdEncoding.DecodeString(encrypted)
	if len(ct) < 16 || string(ct[:8]) != "S_" {
		return ""
	}

	salt := ct[8:16]
	ct = ct[16:]
	key, iv := __DeriveKeyAndIv(passphrase, string(salt))

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	cbc := cipher.NewCBCDecrypter(block, []byte(iv))
	dst := make([]byte, len(ct))
	cbc.CryptBlocks(dst, ct)

	return string(__PKCS7Trimming(dst))
}

func __PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func __PKCS7Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

func __DeriveKeyAndIv(passphrase string, salt string) (string, string) {
	salted := ""
	dI := ""

	for len(salted) < 48 {
		md := md5.New()
		md.Write([]byte(dI + passphrase + salt))
		dM := md.Sum(nil)
		dI = string(dM[:16])
		salted = salted + dI
	}

	key := salted[0:32]
	iv := salted[32:48]

	return key, iv
}

func DecryptJsStr(encryptedStr, keyStr, ivStr string) (string, error) {
	// 将加密字符串从 base64 编码解码
	encryptedData, err := b64.StdEncoding.DecodeString(encryptedStr)
	if err != nil {
		return "", err
	}

	// 确保密钥长度正确
	key := []byte(keyStr)
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", fmt.Errorf("invalid key size %d", len(key))
	}

	// 初始化向量
	iv := []byte(ivStr)
	if len(iv) != aes.BlockSize {
		return "", fmt.Errorf("invalid IV size %d", len(iv))
	}

	// 创建 cipher.Block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 创建 cipher.BlockMode
	mode := cipher.NewCBCDecrypter(block, iv)

	// 解密
	if len(encryptedData)%aes.BlockSize != 0 {
		return "", fmt.Errorf("encrypted data is not a multiple of the block size")
	}
	mode.CryptBlocks(encryptedData, encryptedData)

	// 去除 PKCS7 填充
	padding := encryptedData[len(encryptedData)-1]
	if int(padding) > aes.BlockSize || padding == 0 {
		return "", fmt.Errorf("invalid padding")
	}
	for _, v := range encryptedData[len(encryptedData)-int(padding):] {
		if v != padding {
			return "", fmt.Errorf("invalid padding")
		}
	}
	encryptedData = encryptedData[:len(encryptedData)-int(padding)]

	return string(encryptedData), nil
}
