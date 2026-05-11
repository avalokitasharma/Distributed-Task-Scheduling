package tenantservice

import (
	"log"
	"net/http"
	"os"

	"github.com/avalokitasharma/job-scheduler/common/middleware"
	"github.com/avalokitasharma/job-scheduler/common/postgres"
	"github.com/avalokitasharma/job-scheduler/common/redis"
	"github.com/avalokitasharma/job-scheduler/tenant-service/handlers"
	"github.com/avalokitasharma/job-scheduler/tenant-service/repository"
	"github.com/avalokitasharma/job-scheduler/tenant-service/service"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	db := postgres.ConnectToDB(dsn)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET missing")
	}
	redisAddr := os.Getenv("REDIS_ADDR")
	redisClient, _ := redis.New(redisAddr)

	repo := repository.NewTenantConfigRepo(db)
	svc := service.NewTenantConfigService(repo, redisClient)
	handler := handlers.NewTenantConfigHandler(svc)

	mux := http.NewServeMux()

	// protected routes
	mux.Handle(
		"PUT /tenant/config",
		middleware.AuthMiddleware(secret, http.HandlerFunc(handler.UpsertConfig)),
	)

	mux.Handle(
		"GET /tenant/config",
		middleware.AuthMiddleware(secret, http.HandlerFunc(handler.GetConfig)),
	)

	mux.Handle(
		"GET /tenant/quota/job",
		middleware.AuthMiddleware(secret, http.HandlerFunc(handler.CanCreateJob)),
	)

	mux.Handle(
		"GET /tenant/quota/run",
		middleware.AuthMiddleware(secret, http.HandlerFunc(handler.CanRunJob)),
	)

	mux.Handle(
		"GET /tenant/rate-limit",
		middleware.AuthMiddleware(secret, http.HandlerFunc(handler.CheckRateLimit)),
	)

	log.Println("tenant-service running on :8082")
	err := http.ListenAndServe(":8082", mux)
	if err != nil {
		log.Fatal(err)
	}
}
