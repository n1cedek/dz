package converter

import (
	"dz/auth/internal/model"
	modelRepo "dz/auth/internal/repo/auth/model"
)

func ToAuthFromRepo(auth *modelRepo.User) *model.User {
	return &model.User{
		ID:        auth.ID,
		Name:      auth.Name,
		Email:     auth.Email,
		Password:  auth.Password,
		Role:      ToAuthRoleFromRepo(auth.Role),
		CreatedAt: auth.CreatedAt,
		UpdatedAt: auth.UpdatedAt,
	}
}

func ToAuthRoleFromRepo(role modelRepo.Role) model.Role {
	return model.Role{Value: role.String()}
}

func ToPublicInfo(pi *modelRepo.User) *model.PublicInfo {
	return &model.PublicInfo{
		ID:    pi.ID,
		Name:  pi.Name,
		Email: pi.Email,
	}
}
