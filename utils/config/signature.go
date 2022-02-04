package config

import (
	"github.com/spf13/viper"
)

type signature struct {
	Secret    string
	TimeLimit int
}

var Signature signature

func init() {
	Signature = signature{
		Secret:    viper.GetString("signature.secret"),
		TimeLimit: viper.GetInt("signature.time_limit"),
	}
}
