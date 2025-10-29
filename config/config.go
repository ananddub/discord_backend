package config

import (
	"fmt"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

type PostgreSQLStruct struct {
	URL            string `koanf:"url"`
	MaxConnections int    `koanf:"max_connections"`
	MinConnections int    `koanf:"min_connections"`
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

type Config struct {
	Database DatabaseStruct `koanf:"database"`
	Service  ServiceStruct  `koanf:"service"`
}

var cfg *Config = nil

func Load() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}
	k := koanf.New(".")

	if err := k.Load(file.Provider("config.yml"), yaml.Parser()); err != nil {
		return nil, err
	}
	var cfgl Config
	if err := k.Unmarshal("", &cfgl); err != nil {
		fmt.Println("err %v", err)
		return nil, err
	}
	fmt.Println("your data ", cfgl)
	cfg = &cfgl
	return &cfgl, nil
}
