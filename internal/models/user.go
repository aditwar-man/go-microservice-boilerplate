package models

import (
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User full model
// User model with enhanced fields
type User struct {
	ID          int              `json:"id" db:"id"`
	Username    string           `json:"username" db:"username"`
	Email       string           `json:"email" db:"email"`
	Password    string           `json:"-" db:"password_hash"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at"`
	LoginAt     *time.Time       `json:"login_at,omitempty" db:"login_at"`
	Roles       []Role           `json:"roles,omitempty"`
	Permissions []RolePermission `json:"permissions,omitempty"`
}

// UserRole junction table
type UserRole struct {
	UserID int `json:"user_id" db:"user_id"`
	RoleID int `json:"role_id" db:"role_id"`
}
type UserWithRole struct {
	User User `json:"user" db:"user"`
	Role Role `json:"role" db:"role"`
}

// Hash user password with bcrypt
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Compare user password and payload
func (u *User) ComparePasswords(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

// Sanitize user password
func (u *User) SanitizePassword() {
	u.Password = ""
}

// Prepare user for register
func (u *User) PrepareCreate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Password = strings.TrimSpace(u.Password)

	if err := u.HashPassword(); err != nil {
		return err
	}

	return nil
}

// Prepare user for register
func (u *User) PrepareUpdate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))

	return nil
}

// All Users response
type UsersList struct {
	TotalCount int     `json:"total_count"`
	TotalPages int     `json:"total_pages"`
	Page       int     `json:"page"`
	Size       int     `json:"size"`
	HasMore    bool    `json:"has_more"`
	Users      []*User `json:"users"`
}

// Find user query
type UserWithToken struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}
