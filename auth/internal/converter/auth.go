package converter

import (
	"dz/auth/internal/model"
	modelRepo "dz/auth/internal/repo/auth/model"
	desc "dz/auth/pkg/w1"
)

func ToAuthFromServ(auth *model.User) *desc.CreateRequest {
	return &desc.CreateRequest{
		Name:            auth.Name,
		Email:           auth.Email,
		Password:        auth.Password,
		PasswordConfirm: "",
		Role:            0,
	}
}

func ToRepoUserFromCreateRequest(req *desc.CreateRequest) *modelRepo.User {
	return &modelRepo.User{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Role:     modelRepo.Role{Value: req.GetRole().String()}, // если enum
	}
}
func ToGetResponse(info *model.PublicInfo) *desc.GetResponse {
	if info == nil {
		return nil
	}
	return &desc.GetResponse{
		Id:    info.ID,
		Name:  info.Name,
		Email: info.Email,
	}
}

func ToInfoFromUpdate(info *desc.UpdateRequest) *model.User {
	return &model.User{
		ID:    info.Id,
		Name:  info.Name.GetValue(),
		Email: info.Email.GetValue(),
	}
}
