package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

func (header *Header) Gen() (string, error) {
	var (
		err        error
		headerByte []byte
		headerJson string
	)
	headerByte, err = json.Marshal(header)
	if err != nil {
		return headerJson, err
	}

	headerJson = base64.StdEncoding.EncodeToString(headerByte)

	return headerJson, err
}

type Payload struct {
	Exp  time.Time `form:"exp" json:"exp"`
	Name string    `form:"name" json:"name"`
	Id   uint64    `form:"id" json:"id"`
	// 接收者
	Aud string `form:"aud" json:"aud"`
	// 主题
	Sub string `form:"sub" json:"sub"`
	// 发送者
	Iss string    `form:"iss" json:"iss"`
	Iat time.Time `form:"iat" json:"iat"`
}

func (payload *Payload) Gen() (string, error) {
	var (
		err         error
		payloadByte []byte
		payloadJson string
	)
	payloadByte, err = json.Marshal(payload)
	if err != nil {
		return payloadJson, err
	}

	payloadJson = base64.StdEncoding.EncodeToString(payloadByte)

	return payloadJson, err
}

func HmacSha256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))

	return hex.EncodeToString(h.Sum(nil))
}

func Signature(header string, payload string, hashType string, secret string) (string, error) {

	var (
		builder strings.Builder
		combine string
		err     error
	)

	_, err = builder.WriteString(header)
	if err != nil {
		return "", err
	}
	_, err = builder.WriteString(payload)
	if err != nil {
		return "", err
	}
	combine = builder.String()

	switch hashType {
	case "HS256":
		return HmacSha256(combine, secret), nil
	default:
		return "", errors.New("签名方式不支持")
	}

}
