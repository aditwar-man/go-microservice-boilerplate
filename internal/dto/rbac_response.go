package dto

import (
	"time"

	"github.com/aditwar-man/go-microservice-boilerplate/internal/models"
)

// Check responses
type PermissionCheckResponse struct {
	UserID         int       `json:"user_id"`
	PermissionName string    `json:"permission_name"`
	ResourceName   string    `json:"resource_name"`
	ContextName    *string   `json:"context_name,omitempty"`
	HasPermission  bool      `json:"has_permission"`
	CheckedAt      time.Time `json:"checked_at"`
}

type RoleCheckResponse struct {
	UserID    int       `json:"user_id"`
	RoleName  string    `json:"role_name"`
	HasRole   bool      `json:"has_role"`
	CheckedAt time.Time `json:"checked_at"`
}

// Report responses
type RoleUsageReport struct {
	RoleID         int    `json:"role_id"`
	RoleName       string `json:"role_name"`
	UserCount      int    `json:"user_count"`
	DirectUsers    int    `json:"direct_users"`
	InheritedUsers int    `json:"inherited_users"`
}

type PermissionUsageReport struct {
	PermissionID   int    `json:"permission_id"`
	PermissionName string `json:"permission_name"`
	ResourceName   string `json:"resource_name"`
	RoleCount      int    `json:"role_count"`
	UserCount      int    `json:"user_count"`
}

// System responses
type SystemRolesResponse struct {
	Roles      []models.Role `json:"roles"`
	TotalCount int           `json:"total_count"`
}

type SystemPermissionsResponse struct {
	Permissions []models.Permission `json:"permissions"`
	TotalCount  int                 `json:"total_count"`
}

type UserWithRolesResponse struct {
	User  models.User   `json:"user"`
	Roles []models.Role `json:"roles"`
}

type AllUsersWithRolesResponse struct {
	Users      []UserWithRolesResponse `json:"users"`
	TotalCount int                     `json:"total_count"`
}

// Bulk operation responses
type BulkOperationResponse struct {
	SuccessCount int       `json:"success_count"`
	FailureCount int       `json:"failure_count"`
	Errors       []string  `json:"errors,omitempty"`
	ProcessedAt  time.Time `json:"processed_at"`
}
