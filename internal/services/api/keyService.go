package api

import (
	"cmsApp/configs"
	"cmsApp/internal/constant"
	"cmsApp/internal/dao"
	"cmsApp/pkg/captcha"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"sync"
)

type apiKeyService struct {
	Dao *dao.AppUserDao
}

var (
	instanceApiKeyService *apiKeyService
	onceApiKeyService     sync.Once
)

func NewApiKeyService() *apiKeyService {
	onceApiKeyService.Do(func() {
		instanceApiKeyService = &apiKeyService{
			Dao: dao.NewAppUserDao(),
		}
	})
	return instanceApiKeyService
}

func (ser *apiKeyService) EncryptPasswordPublicKey(password string) (publicKey string, err error) {
	// 解析pem格式的公钥数据
	blockPub, _ := pem.Decode([]byte(configs.App.Rsa.PublicStr))
	if err != nil {
		return "", errors.New(constant.RSA_LOAD_PUBLIC_KEY_ERR)
	}
	pkey, _ := x509.ParsePKCS1PublicKey(blockPub.Bytes)
	if pkey == nil {
		return "", errors.New(constant.RSA_PARSE_PUBLIC_ERR)
	}
	in, err := rsa.EncryptPKCS1v15(rand.Reader, pkey, []byte(password))
	if err != nil {
		return "", err
	}
	return string(in), nil
}

func (ser *apiKeyService) DecryptPasswordPrivateKey(cipherData string) (password string, err error) {
	plainText, err := base64.StdEncoding.DecodeString(cipherData)
	block, _ := pem.Decode([]byte(configs.App.Rsa.PrivateStr))
	if block == nil {
		return "", errors.New(constant.RSA_LOAD_PRIVATE_KEY_ERR)
	}
	// 解析 RSA 私钥
	pkcs8PrivateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	privateKey := pkcs8PrivateKey.(*rsa.PrivateKey)
	if err != nil {
		return "", errors.New(constant.RSA_PARSE_PRIVATE_ERR)
	}
	decryptText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, plainText)
	if err != nil {
		return "", err
	}
	return string(decryptText), nil
}

func (ser *apiKeyService) GenerateSimpleCaptchaCode(id string) (code, b64s, codeId string, err error) {
	code, b64s, codeId, err = captcha.CreateCaptcha(id, "LOGIN", 5)
	return code, b64s, codeId, err
}
