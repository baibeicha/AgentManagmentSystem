package config

import (
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper
	filename   string
	App        AppConfig        `mapstructure:"app"`
	GRPC       GRPCConfig       `mapstructure:"grpc"`
	Datasource DatasourceConfig `mapstructure:"datasource"`
	Redis      RedisConfig      `mapstructure:"redis"`
	Superuser  SuperuserConfig  `mapstructure:"superuser"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Log        LogConfig        `mapstructure:"log"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
}

type GRPCConfig struct {
	Port string    `mapstructure:"port"`
	TLS  TLSConfig `mapstructure:"tls"`
}

type TLSConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	CertPath string `mapstructure:"cert_path"`
	KeyPath  string `mapstructure:"key_path"`
}

type DatasourceConfig struct {
	URL      string `mapstructure:"url"`
	DB       string `mapstructure:"db"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
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

type SuperuserConfig struct {
	Username string `mapstructure:"username"`
	Email    string `mapstructure:"email"`
	Password string `mapstructure:"password"`
}

type JWTConfig struct {
	TTL            JWTTTLConfig `mapstructure:"ttl"`
	PrivateKeyPath string       `mapstructure:"private_key_path"`
	PublicKeyPath  string       `mapstructure:"public_key_path"`
}

type JWTTTLConfig struct {
	Access  int    `mapstructure:"access"`
	Refresh int    `mapstructure:"refresh"`
	Unit    string `mapstructure:"unit"`
}

type LogConfig struct {
	Type  string `mapstructure:"type"`
	Level string `mapstructure:"level"`
	Path  string `mapstructure:"path"`
}

func (t *JWTTTLConfig) GetRefreshTTL() time.Duration {
	return time.Duration(t.Refresh) * GetTimeUnit(t.Unit)
}

func (t *JWTTTLConfig) GetAccessTTL() time.Duration {
	return time.Duration(t.Access) * GetTimeUnit(t.Unit)
}

func init() {
	_ = godotenv.Load()
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
		if defaults[i] == "cfg-path" && i+1 < len(defaults) {
			v.AddConfigPath(defaults[i+1])
			defaults = slices.Delete(defaults, i, i+2)
			i--
		}
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("error finding config file: %v", err)
	}

	configFile := v.ConfigFileUsed()

	configBytes, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("error reading config file %s: %v", configFile, err)
	}

	expandedConfig := os.ExpandEnv(string(configBytes))

	configExt := strings.TrimPrefix(filepath.Ext(configFile), ".")
	v.SetConfigType(configExt)

	if err := v.ReadConfig(strings.NewReader(expandedConfig)); err != nil {
		log.Fatalf("error parsing expanded config: %v", err)
	}

	var cfg Config

	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("error parsing config into struct: %v", err)
	}

	for i := 0; i < len(defaults); i += 2 {
		if !v.IsSet(defaults[i]) && i+1 < len(defaults) {
			v.SetDefault(defaults[i], defaults[i+1])
		}
	}

	cfg.Viper = v
	cfg.filename = filename

	return &cfg
}

func GetTimeUnit(unit string) time.Duration {
	switch unit {
	case "s":
		return time.Second
	case "m":
		return time.Minute
	case "h":
		return time.Hour
	case "d":
		return time.Hour * 24
	default:
		return time.Minute
	}
}
