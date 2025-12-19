package http

import (
	"booking-service/internal/config"
	"booking-service/internal/handlers"
	"booking-service/internal/repo/postgres"
	"booking-service/internal/security"
	"booking-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(db *pgxpool.Pool, cfg config.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	jwtm := security.NewJWTManager(cfg.Auth.JWTSecret)
	repo := postgres.NewBookingRepo()
	svc := service.NewBookingService(db, repo)
	bh := handlers.NewBookingHandlers(jwtm, svc)

	r.POST("/booking/create", bh.CreateAuth)

	return r
}
