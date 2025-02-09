package env

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

const (
	development = "dev"
	staging     = "staging"
	production  = "prod"
)

func GetAppEnv() (string, error) {

	appEnv := os.Getenv("APP_ENV")

	if appEnv == development {
		return development, nil
	} else if appEnv == staging {
		return staging, nil
	} else if appEnv == production {
		return production, nil
	} else {
		return "", errors.New("APP_ENV not found")
	}
}

func LoadConfigFile(appEnv string) error {
	err := godotenv.Load(appEnv + ".env")
	return err
}
