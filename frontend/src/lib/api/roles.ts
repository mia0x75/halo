/** 角色名与后端权限语义保持一致：ADMIN 对应 Root 与 Admin。 */
const adminRoleNames = new Set(["root", "admin"])

export interface RoleRef {
  Name: string
}

export function isAdmin(roles: RoleRef[] | null | undefined): boolean {
  return (roles ?? []).some((role) => adminRoleNames.has(role.Name.trim().toLowerCase()))
}
