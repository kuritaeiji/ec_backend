package config

import (
	"fmt"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/joho/godotenv"
)

func SetupEnv() error {
	env := os.Getenv("ENV")
	wd := os.Getenv("WD")
	err := godotenv.Load(fmt.Sprintf("%s/env/%s.env", wd, env))
	if err != nil {
		return errors.WithStack(err)
	}

	if env == "dev" || env == "test" {
		err = godotenv.Load(fmt.Sprintf("%s/env/secret.env", wd))
	}

	return errors.WithStack(err)
}
