package config

import (
	"log"
	"slices"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper
	filename string
	GRPC     GRPCConfig  `mapstructure:"grpc"`
	Redis    RedisConfig `mapstructure:"redis"`
}

type GRPCConfig struct {
	Port int       `mapstructure:"port"`
	TLS  TLSConfig `mapstructure:"tls"`
}

type TLSConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	CertPath string `mapstructure:"cert_path"`
	KeyPath  string `mapstructure:"key_path"`
}

type RedisConfig struct {
	Addr        string        `mapstructure:"addr"`
	Password    string        `mapstructure:"password"`
	User        string        `mapstructure:"user"`
	DB          int           `mapstructure:"db"`
	MaxRetries  int           `mapstructure:"max_retries"`
	DialTimeout time.Duration `mapstructure:"dial_timeout"`
	Timeout     time.Duration `mapstructure:"timeout"`
}

func MustLoadEmpty(defaults ...string) *Config {
	v := viper.New()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var cfg Config
	cfg.Viper = v
	for i := 0; i < len(defaults); i += 2 {
		if !v.IsSet(defaults[i]) {
			v.SetDefault(defaults[i], defaults[i+1])
		}
	}

	return &cfg
}

func MustLoad(filename string, defaults ...string) *Config {
	if filename == "" {
		filename = "config"
	}

	v := viper.New()

	v.SetConfigName(filename)

	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("./configs")

	for i := 0; i < len(defaults); i++ {
		if defaults[i] == "cfg-path" {
			v.AddConfigPath(defaults[i+1])
			defaults = slices.Delete(defaults, i, i+2)
		}
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	var cfg Config

	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("error parsing config into struct: %v", err)
	}

	for i := 0; i < len(defaults); i += 2 {
		if !v.IsSet(defaults[i]) {
			v.SetDefault(defaults[i], defaults[i+1])
		}
	}

	cfg.Viper = v
	cfg.filename = filename

	return &cfg
}

func (c *Config) GetServerPort() uint16 {
	return uint16(c.GetInt("server.port"))
}
