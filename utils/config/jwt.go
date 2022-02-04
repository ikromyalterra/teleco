package config

import (
	"github.com/spf13/viper"
)

type (
	// JWTConfig JWTConfig
	JWTConfig struct {
		SignKey             []byte
		RefreshSignKey      []byte
		ExpiredToken        int
		ExpiredRefreshToken int
	}
)

func GetJWTSigningKey() []byte {
	return []byte(viper.GetString("jwt.secret_token"))
}

func GetJWTRefreshSigningKey() []byte {
	return []byte(viper.GetString("jwt.secret_refresh_token"))
}

func GetJWTExpiredToken() int {
	return viper.GetInt("jwt.expired_token")
}
func GetJWTExpiredRefreshToken() int {
	return viper.GetInt("jwt.expired_refresh_token")
}
