package config

import (
	"fmt"
	"net/url"
	"os"

	"github.com/jhunt/go-log"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Name    string            `yaml:"name"`
	URL     string            `yaml:"url"`
	Method  string            `yaml:"string,omitempty"`
	Headers map[string]string `yaml:"headers,omitempty"`
	Body    string            `yaml:"body,omitempty"`
}

func NewConfig(filename string) ([]Config, error) {
	// Read from config file
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("No config provided: %s", err)
	}

	// Parse contents of the config file
	var configs []Config
	err = yaml.Unmarshal(bytes, &configs)
	if err != nil {
		return nil, fmt.Errorf("Config not parseable: %s.\nContents are: \n%s", err, string(bytes))
	}

	// Validation on the parsed configs. Ignore bad entries.
	var validateConfigs []Config
	for _, conf := range configs {
		if conf.Name == "" {
			log.Errorf("Missing required name parameter. Skipping health check...")
			continue
		}
		if conf.URL == "" {
			log.Errorf("Missing required url parameter for %s. Skipping health check...", conf.Name)
			continue
		} else {
			_, err := url.Parse(conf.URL)
			if err != nil {
				log.Errorf("Malformed URL %s: %s", conf.URL, err)
				continue
			}
		}
		if conf.Method == "" {
			conf.Method = "GET"
		}

		log.Debugf("Successfully parsed %s config", conf.Name)
		validateConfigs = append(validateConfigs, conf)
	}

	return validateConfigs, nil
}
