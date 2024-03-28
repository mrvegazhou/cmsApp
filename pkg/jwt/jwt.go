package jwt

import (
	"cmsApp/configs"
	"cmsApp/internal/constant"
	"errors"
	"fmt"
	jwtLib "github.com/golang-jwt/jwt/v5"
	"time"
)

// 生成token
func Generate(myClaims MyClaims, secret string) (string, error) {
	if myClaims.IssuedAt == nil {
		myClaims.IssuedAt = jwtLib.NewNumericDate(time.Now())
	}
	if myClaims.NotBefore == nil {
		myClaims.NotBefore = jwtLib.NewNumericDate(time.Now())
	}
	if myClaims.ExpiresAt == nil {
		myClaims.ExpiresAt = jwtLib.NewNumericDate(time.Now().Add(24 * time.Hour))
	}
	if myClaims.Subject == "" {
		myClaims.Subject = "cms"
	}
	if secret == "" {
		secret = configs.App.Login.JwtSecret
	}
	//SetClaims := MyClaims{
	//	Name: name,
	//	//Password: password,
	//	RegisteredClaims: jwt.RegisteredClaims{
	//		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), //有效时间
	//		IssuedAt:  jwt.NewNumericDate(time.Now()),                     //签发时间
	//		NotBefore: jwt.NewNumericDate(time.Now()),                     //生效时间
	//		Issuer:    os.Getenv("JWT_ISSUER"),                            //签发人
	//		Subject:   "somebody",                                         //主题
	//		ID:        "1",                                                //JWT ID用于标识该JWT
	//		Audience:  []string{"somebody_else"},                          //用户
	//	},
	//}

	//使用指定的加密方式和声明类型创建新令牌
	tokenStruct := jwtLib.NewWithClaims(jwtLib.SigningMethodHS256, myClaims)
	tokenStruct.Header["alg"] = "HS256" // 这通常是库自动设置的，但你可以覆盖它
	tokenStruct.Header["typ"] = "JWT"   // JWT 类型

	//获得完整的、签名的令牌
	token, err := tokenStruct.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

/**
* 校验jtoken
 */
func Check(token string, secret string, unlimited bool) (*MyClaims, error) {
	if secret == "" {
		secret = configs.App.Login.JwtSecret
	}
	tokenObj, err := jwtLib.ParseWithClaims(token, &MyClaims{}, func(token *jwtLib.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		if errors.Is(err, jwtLib.ErrTokenMalformed) {
			// 令牌的结构不正确或格式有问题
			return nil, errors.New(constant.TOKEN_MALFORMED_ERR)
		} else if errors.Is(err, jwtLib.ErrTokenNotValidYet) {
			// 令牌没有生效
			return nil, errors.New(constant.TOKEN_NOT_VALID_YET)
		} else if !errors.Is(err, jwtLib.ErrTokenExpired) {
			fmt.Println(err)
			return nil, errors.New(constant.TOKEN_CHECK_ERR)
		}
		if !unlimited {
			if errors.Is(err, jwtLib.ErrTokenExpired) {
				// 过期
				return nil, errors.New(constant.TOKEN_EXPIRE)
			}
		}
	}
	claims, ok := tokenObj.Claims.(*MyClaims)
	if ok && tokenObj.Valid {
		return claims, nil
	} else {
		if !tokenObj.Valid {
			if claims, ok := tokenObj.Claims.(*jwtLib.RegisteredClaims); ok {
				if !unlimited && time.Now().After(claims.ExpiresAt.Time) {
					return nil, errors.New(constant.TOKEN_EXPIRE)
				} else {
					return nil, errors.New(constant.TOKEN_CHECK_ERR)
				}
			}
		}
	}
	return nil, errors.New(constant.TOKEN_CHECK_ERR)
}
