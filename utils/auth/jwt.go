package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sepulsa/teleco/utils/config"
	"github.com/spf13/viper"
)

type (
	JWT interface {
		CreateToken(tokenClaims JWTClaims) (string, error)
		CreateRefreshToken(refreshTokenClaims JWTRefreshClaims) (string, error)
		ParseToken(tokenString string) (interface{}, *JWTClaims, error)
		ParseRefreshToken(refreshTokenString string) (*JWTRefreshClaims, error)
	}
	JWTConf   config.JWTConfig
	JWTClaims struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		Fullname string `json:"fullname"`
		jwt.StandardClaims
	}
	JWTRefreshClaims struct {
		ID string `json:"id"`
		jwt.StandardClaims
	}
)

var (
	ErrInvalidJWT      error = errors.New("invalid or expired token")
	TestJWTAlwaysValid       = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjYzMjg0MDIyMTg2MjkyNTAiLCJlbWFpbCI6InB1dHJpQGFsdGVycmEuaWQiLCJmdWxsbmFtZSI6InB1dHJpIG1hcmlhIiwiZXhwIjoxNjMyODQwMzQxfQ.XWxMkWn5xf2XuiIddTYQaZjRxq2YiopHbR92sG-yvLA"
)

func NewJWT() JWT {
	return &JWTConf{
		SignKey:             config.GetJWTSigningKey(),
		RefreshSignKey:      config.GetJWTRefreshSigningKey(),
		ExpiredToken:        config.GetJWTExpiredToken(),
		ExpiredRefreshToken: config.GetJWTExpiredRefreshToken(),
	}
}

func (j *JWTConf) CreateToken(tokenClaims JWTClaims) (string, error) {
	tokenClaims.ExpiresAt = time.Now().Add(time.Hour * time.Duration(j.ExpiredToken)).Unix()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	return t.SignedString(j.SignKey)
}

func (j *JWTConf) CreateRefreshToken(refreshTokenClaims JWTRefreshClaims) (string, error) {
	refreshTokenClaims.ExpiresAt = time.Now().Add(time.Hour * time.Duration(j.ExpiredRefreshToken)).Unix()
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	return rt.SignedString(j.RefreshSignKey)
}

func (j *JWTConf) ParseToken(tokenString string) (interface{}, *JWTClaims, error) {
	tokenClaims := new(JWTClaims)
	// hack for testing
	if strings.Contains(viper.GetString("env"), "testing") && tokenString == TestJWTAlwaysValid {
		return new(jwt.Token), tokenClaims, nil
	}
	token, err := jwt.ParseWithClaims(tokenString, tokenClaims, func(token *jwt.Token) (interface{}, error) {
		return j.SignKey, nil
	})
	if err != nil || !token.Valid {
		return nil, nil, ErrInvalidJWT
	}

	return token, tokenClaims, nil
}

func (j *JWTConf) ParseRefreshToken(refreshTokenString string) (*JWTRefreshClaims, error) {
	refreshTokenClaims := new(JWTRefreshClaims)
	// hack for testing
	if strings.Contains(viper.GetString("env"), "testing") && refreshTokenString == TestJWTAlwaysValid {
		return refreshTokenClaims, nil
	}
	refreshToken, err := jwt.ParseWithClaims(refreshTokenString, refreshTokenClaims, func(token *jwt.Token) (interface{}, error) {
		return j.RefreshSignKey, nil
	})
	if err != nil || !refreshToken.Valid {
		return nil, ErrInvalidJWT
	}
	return refreshTokenClaims, nil
}
