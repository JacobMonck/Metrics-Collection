package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

func EnvOr(key, or string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return or
}

func RequireEnv(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	logrus.Fatalf("Missing required environment variable: %s", key)
	return ""
}
