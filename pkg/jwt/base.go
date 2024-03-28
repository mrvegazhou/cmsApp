package jwt

import (
	jwtLib "github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	Name string `form:"name" json:"name"`
	jwtLib.RegisteredClaims
}
