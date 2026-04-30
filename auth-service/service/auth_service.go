package service

import (
	"errors"

	"github.com/avalokitasharma/job-scheduler/auth-service/repository"
	"github.com/avalokitasharma/job-scheduler/common/auth"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo   *repository.UserRepo
	tenantRepo *repository.TenantRepo
	secret     string
}

func NewAuthService(userRepo *repository.UserRepo, tenantRepo *repository.TenantRepo, secret string) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		tenantRepo: tenantRepo,
		secret:     secret,
	}
}

func (s *AuthService) RegisterTenant(email, password, tenantName string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	tenant := &repository.Tenant{
		ID:   uuid.NewString(),
		Name: tenantName,
	}

	err = s.tenantRepo.CreateTenant(tenant)
	if err != nil {
		return "", err
	}

	user := &repository.User{
		ID:       uuid.NewString(),
		Email:    email,
		Password: string(hash),
		TenantID: tenant.ID,
		Role:     "admin", // first user = admin
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return "", err
	}

	return auth.GenerateJWT(s.secret, user.ID, tenant.ID, user.Role)
}

func (s *AuthService) RegisterUser(email, password, tenantID, role string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &repository.User{
		ID:       uuid.NewString(),
		Email:    email,
		Password: string(hash),
		TenantID: tenantID,
		Role:     role,
	}

	return s.userRepo.Create(user)
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetByMail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	return auth.GenerateJWT(s.secret, user.ID, user.TenantID, user.Role)
}
