package users

import (
	"github.com/guemidiborhane/factorydigitale.tech/internal/logger"
	"github.com/guemidiborhane/factorydigitale.tech/internal/setup"
	"github.com/guemidiborhane/factorydigitale.tech/pkg/permissions"
	"github.com/guemidiborhane/factorydigitale.tech/pkg/users/auth"
	"gorm.io/gorm"
)

func setupPermissions() {
	permissions.RegisterPermissions("users", []string{"index", "create"})
	permissions.RegisterPermissions("permissions", permissions.DefaultActions)
}

func Setup(c *setup.Config) {
	permissions.SetupCasbin(*c.Database)
	setupPermissions()
	SetupRoutes(*c.Router)
	auth.SetupModels(*c.Database)
	createDefaultUsers(*c.Database)
}

func createUser(username string, password string) error {
	user := auth.User{
		Username: username,
		Password: password,
		Role:     username,
	}

	if err := user.Create(); err != nil {
		return err
	}

	return nil
}

func createDefaultUsers(db *gorm.DB) error {
	var users []auth.User

	if err := db.Select("id").Find(&users).Error; err != nil {
		return err
	}

	if len(users) == 0 {
		logger.Info("Creating default users since none are registered yet", logger.Attrs{})
		createUser("root", "password")
		createUser("admin", "password")
	}

	return nil
}
