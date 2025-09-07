package config

import "os"

type Config struct {
	Port           string
	MaxFileSize    int64
	AllowedFormats []string
	UploadDir      string
}

func New() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		Port:           port,
		MaxFileSize:    10 << 20, //10MB
		AllowedFormats: []string{".jpg", ".jpeg", ".png", ".gif"},
		UploadDir:      "web/static/uploads",
	}

}
