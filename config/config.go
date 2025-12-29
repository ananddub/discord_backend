package config

import (
	"os"
	"path/filepath"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

type PostgreSQLStruct struct {
	URL            string `koanf:"url"`
	MaxConnections int    `koanf:"max_connections"`
	MinConnections int    `koanf:"min_connections"`
}

type ReactiveService struct {
	host string `koanf:"host"`
	port string `koanf:"port"`
}

type RedisStruct struct {
	URL string `koanf:"url"`
}

type DatabaseStruct struct {
	PostgreSQL PostgreSQLStruct `koanf:"postgresql"`
	Redis      RedisStruct      `koanf:"redis"`
}
type ServiceStruct struct {
	Environment string `koanf:"environment"`
	Port        string `koanf:"port"`
}
type S3Struct struct {
	Endpoint  string `koanf:"endpoint"`
	Bucket    string `koanf:"bucket"`
	AccessKey string `koanf:"accessKey"`
	SecretKey string `koanf:"secretKey"`
	UseSSL    bool   `koanf:"useSSL"`
}
type Config struct {
	Database DatabaseStruct  `koanf:"database"`
	Service  ServiceStruct   `koanf:"service"`
	S3       S3Struct        `koanf:"s3"`
	reactive ReactiveService `koanf:"reactive"`
}

var cfg *Config = nil

func Load() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}
	k := koanf.New(".")
	path := GetProjectRoot()
	pathjoin := filepath.Join(path, "config.yml")
	if err := k.Load(file.Provider(pathjoin), yaml.Parser()); err != nil {
		return nil, err
	}
	var cfgl Config
	if err := k.Unmarshal("", &cfgl); err != nil {
		return nil, err
	}
	cfg = &cfgl
	return &cfgl, nil
}
func GetProjectRoot() string {
	dir, _ := os.Getwd()

	for {
		// go.mod exists here? then this is root.
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		// Move one directory up
		parent := filepath.Dir(dir)
		if parent == dir {
			panic("go.mod not found")
		}
		dir = parent
	}
}
