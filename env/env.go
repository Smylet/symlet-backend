package env

import "os"

func GetEnvFileName() string {
	if os.Getenv("ENV") == "test" {
		return "app_test"
	}
	return "app"
}
