package permissions

import (
	"fmt"
	"strings"

	fiberCasbin "github.com/gofiber/contrib/casbin"
	"github.com/gofiber/fiber/v2"
	"github.com/guemidiborhane/factorydigitale.tech/internal/utils"
	"github.com/guemidiborhane/factorydigitale.tech/pkg/users/auth"
)

func CheckPermission(permissions ...string) fiber.Handler {
	return Authorization.Middleware.RequiresPermissions(permissions, fiberCasbin.WithValidationRule(fiberCasbin.MatchAllRule))
}

func AddPolicy(resource string, action string) error {
	if _, err := Authorization.Enforcer.AddPolicy("root", resource, action); err != nil {
		return err
	}

	return nil
}

func GetCurrentUserPermissions(c *fiber.Ctx) []string {
	user, _ := auth.GetCurrentUser(c)
	policies := Authorization.Enforcer.GetFilteredPolicy(0, user.Role)
	var permissions []string

	for _, entry := range policies {
		resource, action := entry[1], entry[2]
		permissions = append(permissions, fmt.Sprintf("%s:%s", resource, action))

	}

	return permissions
}

func GetRolePermissions(role string) []string {
	var permissions []string
	policies := Authorization.Enforcer.GetFilteredPolicy(0, role)

	for _, entry := range policies {
		resource, action := entry[1], entry[2]
		permissions = append(permissions, fmt.Sprintf("%s:%s", resource, action))

	}

	return permissions
}

type RolesPermissionsMap map[string]map[string][]string

func GetAllRolesPermissions() RolesPermissionsMap {
	permissions := make(RolesPermissionsMap)
	policies := Authorization.Enforcer.GetPolicy()

	for _, entry := range policies {
		subject, resource, action := entry[0], entry[1], entry[2]
		if _, ok := permissions[subject]; !ok {
			permissions[subject] = make(map[string][]string)
		}

		if _, ok := permissions[subject][resource]; !ok {
			permissions[subject][resource] = []string{}
		}
		permissions[subject][resource] = append(permissions[subject][resource], action)
	}

	var roles []string
	if err := db.Model(auth.User{}).Select("role").Find(&roles).Error; err != nil {
		utils.WriteToStderr(err)
	}

	for _, role := range roles {
		// check if role is key in permissions
		if _, ok := permissions[role]; !ok {
			permissions[role] = make(map[string][]string)
		}
	}

	return permissions
}

func GetAllPermissions() map[string][]string {
	permissions := GetRolePermissions("root")
	output := make(map[string][]string)
	for _, action := range permissions {
		parts := strings.Split(action, ":")
		output[parts[0]] = append(output[parts[0]], parts[1])
	}

	return output
}

var DefaultActions = []string{"index", "create", "update", "destroy"}

func RegisterPermissions(resource string, actions []string) error {
	for _, action := range actions {
		if err := AddPolicy(resource, action); err != nil {
			return err
		}
	}
	return nil
}
