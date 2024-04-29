package app

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/hoeg/bluepill/internal/morpheus"
)

type Config struct {
	HTTPConfig        *HTTPConfig
	enforcementConfig morpheus.EnforcementConfig
}

type HTTPConfig struct {
	Port     string
	CertFile string
	KeyFile  string
}

var ErrMissingConfig = os.ErrInvalid

func LoadConfig() (*Config, error) {
	port := os.Getenv("BLUEPILL_HTTP_PORT")
	certificateFile := os.Getenv("BLUEPILL_HTTP_CERTIFICATE_FILE")
	keyFile := os.Getenv("BLUEPILL_HTTP_KEY_FILE")

	whitlistFile := os.Getenv("BLUEPILL_ENFORCEMENT_WHITELIST_FILE")
	enforce := os.Getenv("BLUEPILL_ENFORCEMENT_ENFORCE") != ""

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

	if whitlistFile == "" {
		log.Fatalln("BLUEPILL_ENFORCEMENT_WHITELIST_FILE is not set, cannot start server")
		return nil, ErrMissingConfig
	}

	whitelist, err := ReadWhitelistFile(whitlistFile)
	if err != nil {
		log.Fatalf("Failed to read whitelist file: %s", err)
		return nil, err
	}

	return &Config{
		HTTPConfig: &HTTPConfig{
			Port:     port,
			CertFile: certificateFile,
			KeyFile:  keyFile,
		},
		enforcementConfig: morpheus.EnforcementConfig{
			Whitelist: whitelist,
			Enforce:   enforce,
		},
	}, nil
}

func ReadWhitelistFile(file string) ([]string, error) {
	whitelist := make([]string, 0)

	// Open the file
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line into name and IP
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}

		// Add the IP to the whitelist
		whitelist = append(whitelist, parts[1])
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return whitelist, nil
}
