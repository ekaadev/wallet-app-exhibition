package converter

import (
	"backend/internal/entity"
	"backend/internal/model"
)

func UserToUserResponse(user *entity.User, token string) *model.UserResponse {
	return &model.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Token:    token,
	}
}
