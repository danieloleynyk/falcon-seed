package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Load returns Configuration struct
func Load(path string) (*Configuration, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil || len(bytes) == 0 {
		return nil, fmt.Errorf("error reading config file, %s", err)
	}

	var cfg = new(Configuration)

	if err := yaml.Unmarshal(bytes, cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}

	return cfg, nil
}

// Configuration holds data necessary for configuring application
type Configuration struct {
	Logging *Logging     `yaml:"logging,omitempty"`
	Server  *Server      `yaml:"server,omitempty"`
	DB      *Database    `yaml:"database,omitempty"`
	JWT     *JWT         `yaml:"jwt,omitempty"`
	App     *Application `yaml:"application,omitempty"`
}

// Logging holds data necessary for logger configuration
type Logging struct {
	Path string `yaml:"path,omitempty"`
	// Megabytes
	MaxSize int `yaml:"max_size,omitempty"`
	// Files
	MaxBackups int `yaml:"max_backups,omitempty"`
	// Days
	MaxAge int `yaml:"max_age,omitempty"`
}

// Database holds data necessary for database configuration
type Database struct {
	LogQueries bool `yaml:"log_queries,omitempty"`
	Timeout    int  `yaml:"timeout_seconds,omitempty"`
}

// Server holds data necessary for server configuration
type Server struct {
	Port         int  `yaml:"port,omitempty"`
	Debug        bool `yaml:"debug,omitempty"`
	ReadTimeout  int  `yaml:"read_timeout_seconds,omitempty"`
	WriteTimeout int  `yaml:"write_timeout_seconds,omitempty"`
}

// JWT holds data necessary for JWT configuration
type JWT struct {
	AccessToken  *Token `yaml:"access_token,omitempty"`
	RefreshToken *Token `yaml:"refresh_token,omitempty"`
}

// Token holds data necessary for jwt token configuration
type Token struct {
	MinSecretLength  int    `yaml:"min_secret_length,omitempty"`
	DurationMinutes  int    `yaml:"duration_minutes,omitempty"`
	SigningAlgorithm string `yaml:"signing_algorithm,omitempty"`
}

// Application holds application configuration details
type Application struct {
	MinPasswordStr int    `yaml:"min_password_strength,omitempty"`
	SwaggerUIPath  string `yaml:"swagger_ui_path,omitempty"`
}
