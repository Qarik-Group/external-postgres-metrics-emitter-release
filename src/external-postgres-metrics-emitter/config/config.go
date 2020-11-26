package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LoggregatorConfig LoggregatorConfig `yaml:"loggregator"`
	DatabaseConfigs   []DatabaseConfig  `yaml:"databases"`
}

type LoggregatorConfig struct {
	Address string   `yaml:"address"`
	TLS     TLSCerts `yaml:"tls"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type TLSCerts struct {
	KeyFile    string `yaml:"key_file" json:"keyFile"`
	CertFile   string `yaml:"cert_file" json:"certFile"`
	CACertFile string `yaml:"ca_file" json:"caCertFile"`
}

func LoadConfig(path string) (Config, error) {
	var config Config

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}

	return config, yaml.Unmarshal(yamlFile, &config)

}
