package des

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"encoding/hex"
	"errors"
	log "github.com/sirupsen/logrus"
	"runtime"
)

// https://github.com/wumansgy/goEncrypt/blob/be70423b635c30f6f8cc79c67bdcbb858b5f54ac/des/descbc_test.go
const Ivdes = "wumansgy"

/*
*
 1. Group plaintext
    DES CBC mode encryption and decryption, is an 8-byte block encryption
    If the group is not an integer multiple of 8, you need to consider completing the 8 bits2.
*/
func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
}

func PKCS5UnPadding(plainText []byte, blockSize int) ([]byte, error) {
	length := len(plainText)
	number := int(plainText[length-1])
	if number >= length || number > blockSize {
		return nil, errors.New("padding size error please check the secret key or iv")
	}
	return plainText[:length-number], nil
}

func PKCS5Padding(plainText []byte, blockSize int) []byte {
	padding := blockSize - (len(plainText) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	newText := append(plainText, padText...)
	return newText
}

func DesCbcEncrypt(plainText, secretKey, ivDes []byte) (cipherText []byte, err error) {
	if len(secretKey) != 8 {
		return nil, errors.New("a eight-length secret key is required")
	}
	block, err := des.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}
	paddingText := PKCS5Padding(plainText, block.BlockSize())

	var iv []byte
	if len(ivDes) != 0 {
		if len(ivDes) != block.BlockSize() {
			return nil, errors.New("a eight-length ivdes key is required")
		} else {
			iv = ivDes
		}
	} else {
		iv = []byte(Ivdes)
	} // Initialization vector
	blockMode := cipher.NewCBCEncrypter(block, iv)

	cipherText = make([]byte, len(paddingText))
	blockMode.CryptBlocks(cipherText, paddingText)
	return cipherText, nil
}

func DesCbcDecrypt(cipherText, secretKey, ivDes []byte) (plainText []byte, err error) {
	if len(secretKey) != 8 {
		return nil, errors.New("a eight-length secret key is required")
	}
	block, err := des.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case runtime.Error:
				log.Errorf("runtime err=%v,Check that the key or text is correct", err)
			default:
				log.Errorf("error=%v,check the cipherText ", err)
			}
		}
	}()

	var iv []byte
	if len(ivDes) != 0 {
		if len(ivDes) != block.BlockSize() {
			return nil, errors.New("a eight-length ivdes key is required")
		} else {
			iv = ivDes
		}
	} else {
		iv = []byte(Ivdes)
	} // Initialization vector
	blockMode := cipher.NewCBCDecrypter(block, iv)

	plainText = make([]byte, len(cipherText))
	blockMode.CryptBlocks(plainText, cipherText)

	unPaddingText, err := PKCS5UnPadding(plainText, block.BlockSize())
	if err != nil {
		return nil, err
	}
	return unPaddingText, nil
}

func DesCbcEncryptBase64(plainText, secretKey, ivAes []byte) (cipherTextBase64 string, err error) {
	encryBytes, err := DesCbcEncrypt(plainText, secretKey, ivAes)
	return base64.StdEncoding.EncodeToString(encryBytes), err
}

func DesCbcEncryptHex(plainText, secretKey, ivAes []byte) (cipherTextHex string, err error) {
	encryBytes, err := DesCbcEncrypt(plainText, secretKey, ivAes)
	return hex.EncodeToString(encryBytes), err
}

func DesCbcDecryptByBase64(cipherTextBase64 string, secretKey, ivAes []byte) (plainText []byte, err error) {
	plainTextBytes, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return []byte{}, err
	}
	return DesCbcDecrypt(plainTextBytes, secretKey, ivAes)
}

func DesCbcDecryptByHex(cipherTextHex string, secretKey, ivAes []byte) (plainText []byte, err error) {
	plainTextBytes, err := hex.DecodeString(cipherTextHex)
	if err != nil {
		return []byte{}, err
	}
	return DesCbcDecrypt(plainTextBytes, secretKey, ivAes)
}
