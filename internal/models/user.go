package models

import (
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User full model
type User struct {
	ID        int       `json:"id" db:"id" redis:"user_id" validate:"required"`
	Username  string    `json:"username,omitempty" db:"username" redis:"username" validate:"omitempty,lte=60"`
	Email     string    `json:"email,omitempty" db:"email" redis:"email" validate:"omitempty,lte=60,email"`
	Password  string    `json:"password,omitempty" db:"password" redis:"password" validate:"omitempty,required"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at" redis:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at" redis:"updated_at"`
	LoginDate time.Time `json:"login_at" db:"login_at" redis:"login_at"`
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
