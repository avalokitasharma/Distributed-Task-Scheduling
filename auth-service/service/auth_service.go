package service

import "github.com/avalokitasharma/job-scheduler/auth-service/repository"

type AuthService struct {
	repo   *repository.UserRepo
	secret string
}

func NewAuthService(r *repository.UserRepo, secret string) *AuthService {
	return &AuthService{repo: r, secret: secret}
}
