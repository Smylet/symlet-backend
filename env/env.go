package env

import (
	_ "embed"
	"os"
)

const (
	ENV_TEST = "app_test"
	ENV_DEV  = "app"
)

//go:embed app_test.env
var testEnv []byte

//go:embed app.env
var devEnv []byte

func GetEnv() []byte {
	if os.Getenv("ENV") == "test" {
		return testEnv
	}
	return devEnv
}
