package auth1

import (
	"dz/auth/repo"
	"dz/auth/service"
)

type serv struct {
	authRepo repo.AuthRepo
}

func NewServ(authRepo repo.AuthRepo) service.AuthService{
	return &serv{authRepo: authRepo}
}
func NewMockService(deps ...interface{}) service.AuthService {
	srv := serv{}

	for _, v := range deps {
		switch s := v.(type) {
		case repo.AuthRepo:
			srv.authRepo = s
		}
	}
	return &srv
}
