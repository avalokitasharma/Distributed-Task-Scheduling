package authservice

import (
	"net/http"
	"os"

	"github.com/avalokitasharma/job-scheduler/auth-service/handlers"
	"github.com/avalokitasharma/job-scheduler/auth-service/repository"
	"github.com/avalokitasharma/job-scheduler/auth-service/service"
	"github.com/avalokitasharma/job-scheduler/common/middleware"
	"github.com/avalokitasharma/job-scheduler/common/postgres"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	secret := os.Getenv("JWT_SECRET")

	db := postgres.ConnectToDB(dsn)

	tenantRepo := repository.NewTenantRepo(db)
	userRepo := repository.NewUserRepo(db)

	svc := service.NewAuthService(userRepo, tenantRepo, secret)

	handler := handlers.NewAuthHandler(svc)

	// Routes
	http.HandleFunc("/auth/register", handler.RegisterTenant)
	http.HandleFunc("/auth/login", handler.Login)
	http.Handle("/auth/users", middleware.AuthMiddleware(secret, http.HandlerFunc(handler.CreateUser)))

}
