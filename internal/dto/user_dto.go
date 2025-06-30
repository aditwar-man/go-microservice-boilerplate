package dto

import (
	"github.com/aditwar-man/go-microservice-boilerplate/internal/models"
)

// Request/Response models
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	RoleIDs  []int  `json:"role_ids,omitempty"`
}

type UpdateUserRequest struct {
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	RoleIDs  []int   `json:"role_ids,omitempty"`
}

type AssignRoleRequest struct {
	UserID  int   `json:"user_id" binding:"required"`
	RoleIDs []int `json:"role_ids" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password_hash" validate:"required"`
	Email    string `json:"email" validate:"omitempty,lte=60,email"`
}

type LoginUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token       string                  `json:"token"`
	User        models.User             `json:"user"`
	Permissions []models.RolePermission `json:"permissions"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// Enhanced RBAC Context with hierarchical roles
type RBACContext struct {
	UserID                int                                `json:"user_id"`
	DirectRoles           []models.Role                      `json:"direct_roles"`
	InheritedRoles        []models.Role                      `json:"inherited_roles"`
	AllRoles              []models.Role                      `json:"all_roles"`
	Permissions           []models.RolePermission            `json:"permissions"`
	PermissionsByResource map[string][]models.RolePermission `json:"permissions_by_resource"`
}

// Permission check request
type PermissionCheckRequest struct {
	UserID         int     `json:"user_id"`
	PermissionName string  `json:"permission_name"`
	ResourceName   string  `json:"resource_name"`
	ContextName    *string `json:"context_name,omitempty"`
}

type ContextsList struct {
	TotalCount int              `json:"total_count"`
	TotalPages int              `json:"total_pages"`
	Page       int              `json:"page"`
	Size       int              `json:"size"`
	HasMore    bool             `json:"has_more"`
	Contexts   []models.Context `json:"contexts"`
}

type ResourcesList struct {
	TotalCount int               `json:"total_count"`
	TotalPages int               `json:"total_pages"`
	Page       int               `json:"page"`
	Size       int               `json:"size"`
	HasMore    bool              `json:"has_more"`
	Resources  []models.Resource `json:"resources"`
}

type PermissionsList struct {
	TotalCount  int                 `json:"total_count"`
	TotalPages  int                 `json:"total_pages"`
	Page        int                 `json:"page"`
	Size        int                 `json:"size"`
	HasMore     bool                `json:"has_more"`
	Permissions []models.Permission `json:"permissions"`
}
