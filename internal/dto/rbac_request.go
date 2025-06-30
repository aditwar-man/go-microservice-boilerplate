package dto

// Role management requests
type CreateRoleRequest struct {
	Name         string `json:"name" validate:"required,min=3,max=50"`
	Description  string `json:"description" validate:"max=255"`
	ParentRoleID *int   `json:"parent_role_id,omitempty"`
}

type UpdateRoleRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=3,max=50"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=255"`
}

// Permission management requests
type CreatePermissionRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"max=255"`
}

// Resource management requests
type CreateResourceRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"max=255"`
}

// Context management requests
type CreateContextRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"max=255"`
}

// Role assignment requests
type AssignRolesRequest struct {
	RoleIDs []int `json:"role_ids" validate:"required,min=1"`
}

type AssignRolePermissionRequest struct {
	RoleID       int  `json:"role_id" validate:"required"`
	PermissionID int  `json:"permission_id" validate:"required"`
	ResourceID   int  `json:"resource_id" validate:"required"`
	ContextID    *int `json:"context_id,omitempty"`
}

// Check requests
type CheckPermissionRequest struct {
	UserID         int     `json:"user_id" validate:"required"`
	PermissionName string  `json:"permission_name" validate:"required"`
	ResourceName   string  `json:"resource_name" validate:"required"`
	ContextName    *string `json:"context_name,omitempty"`
}

type CheckRoleRequest struct {
	UserID   int    `json:"user_id" validate:"required"`
	RoleName string `json:"role_name" validate:"required"`
}

// Bulk operation requests
type BulkAssignRolesRequest struct {
	UserIDs []int `json:"user_ids" validate:"required,min=1"`
	RoleIDs []int `json:"role_ids" validate:"required,min=1"`
}

type BulkAssignPermissionsRequest struct {
	RoleIDs      []int `json:"role_ids" validate:"required,min=1"`
	PermissionID int   `json:"permission_id" validate:"required"`
	ResourceID   int   `json:"resource_id" validate:"required"`
	ContextID    *int  `json:"context_id,omitempty"`
}
