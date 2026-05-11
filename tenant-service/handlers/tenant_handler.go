package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/avalokitasharma/job-scheduler/common/middleware"
	"github.com/avalokitasharma/job-scheduler/tenant-service/repository"
	"github.com/avalokitasharma/job-scheduler/tenant-service/service"
)

type TenantConfigHandler struct {
	svc *service.TenantConfigService
}

func NewTenantConfigHandler(s *service.TenantConfigService) *TenantConfigHandler {
	return &TenantConfigHandler{
		svc: s,
	}
}

// Get Config route
func (h *TenantConfigHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	t, err := h.svc.GetConfig(r.Context(), claims.TenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

// PUT config
func (h *TenantConfigHandler) UpsertConfig(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req repository.TenantConfig
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	req.TenantId = claims.TenantID

	err := h.svc.UpsertConfig(r.Context(), claims.TenantID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("tenant config updated"))
}

// GET tenant quota
func (h *TenantConfigHandler) CanCreateJob(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	err := h.svc.CanCreateJob(r.Context(), claims.TenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusTooManyRequests)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("allowed"))
}

// GET check concurrent jobs
func (h *TenantConfigHandler) CanRunJob(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	err := h.svc.CanRunJob(r.Context(), claims.TenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusTooManyRequests)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("allowed"))
}

// GET check rate limit
func (h *TenantConfigHandler) CheckRateLimit(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err := h.svc.CheckRateLimit(r.Context(), claims.TenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusTooManyRequests)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("allowed"))
}
