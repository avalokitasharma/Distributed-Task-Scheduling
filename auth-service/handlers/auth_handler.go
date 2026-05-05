package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/avalokitasharma/job-scheduler/auth-service/service"
	"github.com/avalokitasharma/job-scheduler/common/middleware"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: s}
}

// Register tenant route
type registerTenantReq struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	TenantName string `json:"tenant_name"`
}

func (h *AuthHandler) RegisterTenant(w http.ResponseWriter, r *http.Request) {
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

// Register user route
type createUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	TenantID string `json:"tenant_id"`
	Role     string `json:"role"`
}

func (h *AuthHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	// extract claims from context - set by AuthMiddleware
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	// enforce tenant isolation
	if claims.TenantID != req.TenantID {
		http.Error(w, "forbidden: tenant mismatch", http.StatusForbidden)
		return
	}

	// prevent creating admins
	if req.Role == "admin" {
		http.Error(w, "forbidden: cannot create admin user", http.StatusForbidden)
		return
	}

	err := h.svc.RegisterUser(req.Email, req.Password, req.TenantID, req.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Login route
type loginReq struct {
	Email    string `json:"email" `
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	token, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
