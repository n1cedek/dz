package auth2

import (
	"dz/auth/internal/service"
	desc "dz/auth/pkg/w1"
)

type Implementation struct {
	desc.UnimplementedUserAPIServer
	authService service.AuthService
}

func NewImplementation(noteService service.AuthService) *Implementation {
	return &Implementation{
		authService: noteService,
	}
}
