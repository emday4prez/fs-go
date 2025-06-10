package config

import (
	"os"
)

type Config struct {
	Port      string
	UploadDir string
}

func Load() *Config {
	cfg := &Config{
		Port:      ":8080",
		UploadDir: "./uploads",
	}

	if port, ok := os.LookupEnv("PORT"); ok {
		cfg.Port = ":" + port
	}

	if uploadDir, ok := os.LookupEnv("UPLOAD_DIR"); ok {
		cfg.UploadDir = uploadDir
	}

	return cfg
}
