package env

import (
	"os"
)


func GetEnv(variable string, defaultValue string) string {
	val, ok := os.LookupEnv(variable)

	if ok {
		return val
	}

	return defaultValue
}
