package models

// Role model

// Role model with hierarchical support
type Role struct {
	ID           int              `json:"id" db:"id"`
	Name         string           `json:"name" db:"name"`
	Description  string           `json:"description" db:"description"`
	ParentRoleID *int             `json:"parent_role_id,omitempty" db:"parent_role_id"`
	ParentRole   *Role            `json:"parent_role,omitempty"`
	ChildRoles   []Role           `json:"child_roles,omitempty"`
	Permissions  []RolePermission `json:"permissions,omitempty"`
}

// Enhanced RolePermission with resource and context
type RolePermission struct {
	RoleID       int         `json:"role_id" db:"role_id"`
	PermissionID int         `json:"permission_id" db:"permission_id"`
	ResourceID   int         `json:"resource_id" db:"resource_id"`
	ContextID    *int        `json:"context_id,omitempty" db:"context_id"`
	Role         *Role       `json:"role,omitempty"`
	Permission   *Permission `json:"permission,omitempty"`
	Resource     *Resource   `json:"resource,omitempty"`
	Context      *Context    `json:"context,omitempty"`
}

type RolesList struct {
	TotalCount int     `json:"total_count"`
	TotalPages int     `json:"total_pages"`
	Page       int     `json:"page"`
	Size       int     `json:"size"`
	HasMore    bool    `json:"has_more"`
	Roles      []*Role `json:"roles"`
}
