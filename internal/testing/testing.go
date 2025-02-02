package testing

import (
	"log"
	"testing"

	"github.com/guemidiborhane/factorydigitale.tech/internal/config"
	logger "github.com/guemidiborhane/factorydigitale.tech/internal/logger"
	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

func SetupApp(models ...interface{}) (*fiber.App, *gorm.DB, error) {
	logger.Setup()

	var a *fiber.App
	var db *gorm.DB
	var err error

	return a, db, err
}

func SetupSuite(t *testing.T, models ...interface{}) (*fiber.App, *gorm.DB) {
	a, err := NewApp()
	if err != nil {
		log.Fatal(err)
	}
	if !config.AppConfig.IsTest() {
		log.Fatal("CAUTION: App is not running in test mode")
	}

	db, err := SetupDB(models...)
	if err != nil {
		log.Fatal(err)
	}

	AfterEach(func() {
		CleanDB(models...)
	})
	if err != nil {
		log.Fatal(err)
	}
	RegisterFailHandler(Fail)

	return a, db
}
