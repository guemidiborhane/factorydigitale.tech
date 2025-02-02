package testing

import (
	"fmt"

	"github.com/guemidiborhane/factorydigitale.tech/internal/storage"
	"github.com/bluele/factory-go/factory"
	"gorm.io/gorm"
)

var db *gorm.DB

func SetupDB(models ...interface{}) (*gorm.DB, error) {
	storage.SetupPostgres()
	db = storage.DB

	db.AutoMigrate(models...)

	return db, nil
}

func CleanDB(models ...interface{}) error {
	for _, model := range models {
		stmt := &gorm.Statement{DB: db}
		stmt.Parse(model)
		query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", stmt.Schema.Table)
		if err := db.Exec(query).Error; err != nil {
			return err
		}
	}

	return nil
}

func FactoryPersist(args factory.Args) error {
	return db.Create(args.Instance()).Error
}

func CreateFactory(model interface{}) *factory.Factory {
	return factory.NewFactory(
		model,
	).SeqInt("ID", func(n int) (interface{}, error) {
		return n, nil
	}).OnCreate(FactoryPersist)
}
