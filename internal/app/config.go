package app

import (
	"log"
	"os"
)

type HTTPConfig struct {
	Port     string
	CertFile string
	KeyFile  string
}

var ErrMissingConfig = os.ErrInvalid

func LoadHTTPConfig() (*HTTPConfig, error) {
	port := os.Getenv("BLUEPILL_HTTP_PORT")
	certificateFile := os.Getenv("BLUEPILL_HTTP_CERTIFICATE_FILE")
	keyFile := os.Getenv("BLUEPILL_HTTP_KEY_FILE")

	if port == "" {
		log.Println("BLUEPILL_HTTP_PORT is not set, using default port 8443")
		port = "8443"
	}

	if certificateFile == "" {
		log.Fatalln("BLUEPILL_HTTP_CERTIFICATE_FILE is not set, cannot start server")
		return nil, ErrMissingConfig
	}

	if keyFile == "" {
		log.Fatalln("BLUEPILL_HTTP_KEY_FILE is not set, cannot start server")
		return nil, ErrMissingConfig
	}

	return &HTTPConfig{
		Port:     port,
		CertFile: certificateFile,
		KeyFile:  keyFile,
	}, nil
}
