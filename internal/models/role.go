package models

import (
	"database/sql"
)

type Role struct {
	ID           int           `json:"id" db:"id" redis:"role_id" validate:"required"`
	Name         string        `json:"name" db:"name" redis:"role_name" validate:"omitempty,lte=30"`
	Description  string        `json:"description" db:"description" redis:"description"`
	ParentRoleId sql.NullInt64 `json:"parent_role_id" db:"parent_role_id" redis:"parent_role_id"`
}
