package permissions

import (
	"github.com/guemidiborhane/factorydigitale.tech/internal/setup"
	"gorm.io/gorm"
)

var db *gorm.DB

func Setup(c *setup.Config) {
	db = *c.Database
}
