package config

import (
	"fmt"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/joho/godotenv"
)

func SetupEnv() error {
	env := os.Getenv("ENV")
	err := godotenv.Load(fmt.Sprintf("/go/app/env/%s.env", env))
	return errors.WithStack(err)
}
