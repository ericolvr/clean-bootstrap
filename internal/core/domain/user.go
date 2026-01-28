package domain

import (
	"context"
	"errors"
)

type UserType int

const (
	UserTypeMaster UserType = 1 // Master
	UserTypeAdmin  UserType = 1 // Admin
	UserTypeUser   UserType = 2 // Mobile User

)

type User struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	Mobile   string   `json:"mobile,omitempty"`
	UserType UserType `json:"user_type"`
	Password string   `json:"password,omitempty"`
	Status   bool     `json:"status"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	List(ctx context.Context, limit, offset int) ([]*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByMobile(ctx context.Context, mobile string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
}

func (u *User) ValidateUser() error {
	if u.Name == "" {
		return errors.New("name is required")
	}
	if u.Mobile == "" {
		return errors.New("mobile is required")
	}
	if u.UserType <= 0 {
		return errors.New("invalid user type")
	}
	if u.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

func (u *User) ValidateUserUpdate() error {
	if u.Name == "" {
		return errors.New("name is required")
	}
	if u.Mobile == "" {
		return errors.New("mobile is required")
	}
	if u.UserType <= 0 {
		return errors.New("invalid user type")
	}

	return nil
}
