package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/avalokitasharma/job-scheduler/auth-service/service"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: s}
}

type registerTenantReq struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	TenantName string `json:"tenant_name"`
}

func (h *AuthHandler) RegisterTenant(w http.ResponseWriter, r http.Request) {
	var req registerTenantReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	token, err := h.svc.RegisterTenant(req.Email, req.Password, req.TenantName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

type createUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	TenantID string `json:"tenant_id"`
	Role     string `json:"role"`
}
