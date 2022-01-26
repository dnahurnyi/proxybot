package opts

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func Load() error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get current directory path: %w", err)
	}
	envPath := filepath.Join(dir, ".env")
	if _, err := os.Stat(envPath); err == nil {
		if err := godotenv.Load(envPath); err != nil {
			return fmt.Errorf("load config %w", err)
		}
	}
	return nil
}
