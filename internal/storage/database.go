package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/guemidiborhane/factorydigitale.tech/internal/config"
	"github.com/guemidiborhane/factorydigitale.tech/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type dbConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     uint   `mapstructure:"DB_PORT"`
	Username string `mapstructure:"DB_USERNAME"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
}

var (
	DbConfig = &dbConfig{
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		Password: "postgres",
		Name:     "postgres",
	}
	DB *gorm.DB
)

func SetupPostgres() {
	if err := config.EnvFile.LoadConfig(&DbConfig); err != nil {
		log.Fatal(err)
	}

	if config.AppConfig.IsTest() {
		DbConfig.Name = fmt.Sprintf("%v_test", DbConfig.Name)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
		DbConfig.Host, DbConfig.Username, DbConfig.Password, DbConfig.Name, DbConfig.Port, utils.GetTimeZone(),
	)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		QueryFields: true,
	})

	if err != nil {
		utils.WriteToStderr(err)
	}

	// threads_count := runtime.GOMAXPROCS(0)

	sqlDb, _ := DB.DB()
	sqlDb.SetMaxIdleConns(100)
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetConnMaxIdleTime(10 * time.Minute)

	// DB.Use(
	// 	dbresolver.Register(dbresolver.Config{}).
	// 		SetConnMaxLifetime(10 * time.Minute).
	// 		SetMaxIdleConns(100).
	// 		SetMaxOpenConns(100),
	// )
}
