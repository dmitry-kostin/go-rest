package models

type UserRepository interface {
	CreateUser(user *CreateUserReq) (*User, error)
	ListUsers() (*[]User, error)
	GetUser(userId int64) (*User, error)
	RemoveUser(userId int64) error
}
