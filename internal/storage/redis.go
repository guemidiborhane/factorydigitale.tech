package storage

import (
	"fmt"
	"slices"
	"strings"

	"github.com/guemidiborhane/factorydigitale.tech/internal/config"
	"github.com/guemidiborhane/factorydigitale.tech/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/redis"
	"github.com/spf13/viper"
)

type redisConfig struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     int    `mapstructure:"REDIS_PORT"`
	Database int
}

var (
	RedisConfig = &redisConfig{
		Host: "localhost",
		Port: 6379,
	}
	Redis     *redis.Storage
	Databases = []string{"storage", "sessions", "csrf", "websocket"}
)

func SetupRedis() {
	if err := config.EnvFile.LoadConfig(&RedisConfig); err != nil {
		utils.WriteToStderr(err)
	}

	Redis = RedisStorage("storage")
}

func RedisStorage(database string) *redis.Storage {
	db := slices.Index(Databases, database)

	return redis.New(redis.Config{
		Host:     RedisConfig.Host,
		Port:     RedisConfig.Port,
		Database: db,
	})
}

func WriteRedisEnvVar() {
	if !config.AppConfig.IsDev() || fiber.IsChild() {
		return
	}

	var hosts []string
	for database, dbName := range Databases {
		hosts = append(hosts, fmt.Sprintf("%s:%s:%d:%d", dbName, RedisConfig.Host, RedisConfig.Port, database))
	}

	viper.Set("REDIS_DBS", strings.Join(hosts, ","))
	if err := viper.WriteConfig(); err != nil {
		utils.WriteToStderr(err)
	}
}
