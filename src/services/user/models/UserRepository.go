package models

type UserRepository interface {
	CreateUser(user *User) error
	ListUsers() (*[]User, error)
	GetUser(userId EntityId) (*User, error)
	RemoveUser(userId EntityId) error
}
