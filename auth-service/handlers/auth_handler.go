package handlers

import "github.com/avalokitasharma/job-scheduler/auth-service/service"

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: s}
}
