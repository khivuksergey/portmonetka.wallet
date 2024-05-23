package config

import (
	"fmt"
	"github.com/khivuksergey/portmonetka.wallet/error"
	"github.com/spf13/viper"
)

func LoadEnv() {
	var errMsg error.ErrorMessage

	requiredEnvVars := []string{
		"JWT_SECRET",
		"JWT_ISSUER",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"DB_HOST",
	}

	viper.AutomaticEnv()

	for _, env := range requiredEnvVars {
		if !viper.IsSet(env) {
			errMsg.Append(envErrorMsg(env))
		}
	}

	err := errMsg.ToError()
	if err != nil {
		fmt.Printf("error loading environment variables: %v\n", err)
		panic(err)
	}
}

func envErrorMsg(env string) string { return fmt.Sprintf("%s missing", env) }
