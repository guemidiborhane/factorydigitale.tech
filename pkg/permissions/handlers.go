package permissions

import (
	"slices"
	"strings"

	"github.com/guemidiborhane/factorydigitale.tech/internal/errors"
	"github.com/gofiber/fiber/v2"
)

type PermissionsResponse struct {
	Permissions      map[string][]string `json:"permissions" binding:"required"`
	RolesPermissions RolesPermissionsMap `json:"roles_permissions" binding:"required"`
}

// @Summary	Index Permissions
// @Tags		Permissions
// @Accept		json
// @Produce	json
// @Success	200	{object}	PermissionsResponse
// @Failure	403	{object}	errors.HttpError
// @Failure	500	{object}	errors.HttpError
// @Router		/api/permissions [get]
func Index(c *fiber.Ctx) error {
	permissions := GetAllPermissions()
	rolesPermissions := GetAllRolesPermissions()

	return c.Status(fiber.StatusOK).JSON(&PermissionsResponse{
		Permissions:      permissions,
		RolesPermissions: rolesPermissions,
	})
}

type PermissionsParams map[string][]string

// @Summary	Store Role
// @Tags		Permissions
// @Accept		json
// @Produce	json
// @Param		body	body		PermissionsParams	true	"RoleDefinition"
// @Param		role	path		string				true	"Role"
// @Success	200		{object}	[]string
// @Failure	403		{object}	errors.HttpError
// @Failure	500		{object}	errors.HttpError
// @Router		/api/permissions/{role} [post]
// @Router		/api/permissions/{role} [put]
func Store(c *fiber.Ctx) error {
	var params PermissionsParams
	role := c.Params("role")
	permissions := GetRolePermissions(role)

	if err := c.BodyParser(&params); err != nil {
		return errors.Unexpected(err)
	}

	var toBeRemoved []string = permissions

	for resource, actions := range params {
		for _, action := range actions {
			if action != "" {
				_, err := Authorization.Enforcer.AddPolicy(role, resource, action)
				if err != nil {
					return errors.Unexpected(err)
				}
				if index := slices.Index(toBeRemoved, resource+":"+action); index != -1 {
					toBeRemoved = slices.Delete(toBeRemoved, index, index+1)
				}
			}
		}
	}

	for _, action := range toBeRemoved {
		split := strings.Split(action, ":")
		_, err := Authorization.Enforcer.RemovePolicy(role, split[0], split[1])
		if err != nil {
			return errors.Unexpected(err)
		}
	}

	return c.Status(fiber.StatusOK).JSON(GetRolePermissions(role))
}
