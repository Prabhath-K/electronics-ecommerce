package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("failed to load environment file: %w", err)
	}

	return nil
}

func getEnv(key string) string {
	return os.Getenv(key)
}
