package config

import "os"

func GetConfig(key string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}

	return ""
}
