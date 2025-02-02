package permissions

import (
	"fmt"

	casbin "github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	fiberCasbin "github.com/gofiber/contrib/casbin"
	"github.com/gofiber/fiber/v2"
	"github.com/guemidiborhane/factorydigitale.tech/internal/config"
	"github.com/guemidiborhane/factorydigitale.tech/internal/errors"
	"github.com/guemidiborhane/factorydigitale.tech/pkg/users/auth"
	"gorm.io/gorm"
)

type AuthorizationConfig struct {
	Enforcer   *casbin.Enforcer
	Middleware *fiberCasbin.Middleware
}

var Authorization *AuthorizationConfig

var modelString string = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act || r.sub == "root"
`

func SetupCasbin(db *gorm.DB) {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize casbin adapter: %v", err))
	}

	m, _ := model.NewModelFromString(modelString)

	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		panic(err)
	}

	enforcer.EnableLog(config.AppConfig.IsDev())

	middleware := fiberCasbin.New(fiberCasbin.Config{
		Enforcer: enforcer,

		Lookup: func(c *fiber.Ctx) string {
			user, err := auth.GetCurrentUser(c)
			if err != nil {
				return err.Error()
			}

			return user.Role
		},
		Unauthorized: func(c *fiber.Ctx) error {
			return errors.Unauthorized
		},
		Forbidden: func(c *fiber.Ctx) error {
			return errors.Forbidden
		},
	})

	Authorization = &AuthorizationConfig{
		Enforcer:   enforcer,
		Middleware: middleware,
	}
}
