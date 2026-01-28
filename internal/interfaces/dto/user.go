package dto

import "github.com/ericolvr/sec-backend/internal/core/domain"

type UserRequest struct {
	Name     string          `json:"name" binding:"required"`
	Mobile   string          `json:"mobile" binding:"required"`
	Password string          `json:"password"`
	UserType domain.UserType `json:"user_type" binding:"required,min=1,max=4"`
	Status   bool            `json:"status"`
}

type UserResponse struct {
	ID       int64           `json:"id"`
	Name     string          `json:"name"`
	Mobile   string          `json:"mobile"`
	UserType domain.UserType `json:"user_type"`
	Status   bool            `json:"status"`
}

type UserUpdate struct {
	ID       int64           `json:"id"`
	Name     string          `json:"name" binding:"required"`
	Mobile   string          `json:"mobile" binding:"required"`
	UserType domain.UserType `json:"user_type" binding:"required,min=1,max=4"`
	Status   bool            `json:"status"`
}

func ToUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Mobile:   user.Mobile,
		UserType: user.UserType,
		Status:   user.Status,
	}
}
