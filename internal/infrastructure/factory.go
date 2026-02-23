package infrastructure

import (
    "fmt"

    "github.com/cureerel/gotemplate/internal/infrastructure/dbtypes"
    "github.com/cureerel/gotemplate/internal/infrastructure/postgres"
    redisclient "github.com/cureerel/gotemplate/internal/infrastructure/redis"
)

type DBConfig struct {
    Driver string
    DSN    string
}

type RedisConfig struct {
    Addr     string
    Username string
    Password string
    DB       int
}

func NewDatabase(cfg DBConfig) (dbtypes.DBClient, error) {
    switch cfg.Driver {
    case "postgres":
        return postgres.New(cfg.DSN)
    default:
        return nil, fmt.Errorf("unknown driver: %s", cfg.Driver)
    }
}

func NewRedis(cfg RedisConfig) (redisclient.Redis, error) {
    return redisclient.New(cfg.Addr, cfg.Username, cfg.Password, cfg.DB)
}