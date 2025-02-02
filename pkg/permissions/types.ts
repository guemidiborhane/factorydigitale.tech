import { APIResponse } from "~/api"

export type RolesPermissions = Permissions['roles_permissions']
export type Permissions = APIResponse<'/api/permissions', 'get'>
