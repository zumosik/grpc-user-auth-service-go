package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/zumosik/grpc-user-auth-service-go/storage/postgres"
	"log"
)

type Config struct {
	Postgres postgres.ConnectionData `yaml:"postgres" env-required:"true"` // all fields for db connection
	Addr     string                  `yaml:"addr" env-required:"true"`
	Timeout  int                     `yaml:"timeout" env-required:"true"`
}

// MustConfig will throw fatal errors
func MustConfig(cfgPath string) (cfg Config) {
	err := cleanenv.ReadConfig(cfgPath, &cfg)
	if err != nil {
		log.Fatalf("Cant't read config: %v", err)
	}
	return
}
