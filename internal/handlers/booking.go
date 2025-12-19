package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"booking-service/internal/dto"
	"booking-service/internal/security"
	"booking-service/internal/service"

	"github.com/gin-gonic/gin"
)

type BookingHandlers struct {
	jwt *security.JWTManager
	svc *service.BookingService
}

func NewBookingHandlers(jwt *security.JWTManager, svc *service.BookingService) *BookingHandlers {
	return &BookingHandlers{jwt: jwt, svc: svc}
}

func (h *BookingHandlers) CreateAuth(c *gin.Context) {
	userID, ok := h.userIDFromAuthHeader(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.AuthBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.svc.CreateAuthBooking(ctx, userID, req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRequest) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *BookingHandlers) userIDFromAuthHeader(c *gin.Context) (string, bool) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		return "", false
	}
	const pfx = "Bearer "
	if !strings.HasPrefix(auth, pfx) {
		return "", false
	}
	token := strings.TrimSpace(strings.TrimPrefix(auth, pfx))
	claims, err := h.jwt.Parse(token)
	if err != nil || claims.Subject == "" {
		return "", false
	}
	return claims.Subject, true
}
