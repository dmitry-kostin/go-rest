package models

import (
	"time"
)

type UserRole string

const (
	Admin    UserRole = "admin"
	Customer UserRole = "customer"
)

type User struct {
	Id        int64     `json:"id"`
	Email     string    `json:"email"`
	Role      UserRole  `json:"role"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateUserReq struct {
	Email     string   `json:"email" valid:"email"`
	Role      UserRole `json:"role" valid:"required"`
	FirstName string   `json:"first_name" valid:"required"`
	LastName  string   `json:"last_name" valid:"required"`
}
