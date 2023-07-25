package models

import (
	"github.com/google/uuid"
	"time"
)

type UserRole string

const (
	Admin    UserRole = "admin"
	Customer UserRole = "customer"
)

type EntityId = uuid.UUID

type User struct {
	Id         EntityId  `json:"id"`
	IdentityId uuid.UUID `json:"identity_id"`
	Email      string    `json:"email"`
	Role       UserRole  `json:"role"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewUser(identityId uuid.UUID, email, firstName, lastName string, role UserRole) *User {
	return &User{
		Id:         uuid.New(),
		IdentityId: identityId,
		Email:      email,
		Role:       role,
		FirstName:  firstName,
		LastName:   lastName,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
}
