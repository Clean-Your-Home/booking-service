package service

import (
	"booking-service/internal/dto"
	"booking-service/internal/repo/postgres"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingService struct {
	db   *pgxpool.Pool
	repo *postgres.BookingRepo
}

func NewBookingService(db *pgxpool.Pool, repo *postgres.BookingRepo) *BookingService {
	return &BookingService{db: db, repo: repo}
}

func (s *BookingService) CreateAuthBooking(ctx context.Context, userID string, req dto.AuthBookingRequest) (dto.BookingResponse, error) {
	scheduled, err := time.Parse(time.RFC3339, req.ScheduleInfo.Date)
	if err != nil {
		return dto.BookingResponse{}, fmt.Errorf("%w: bad date", ErrInvalidRequest)
	}
	if scheduled.Before(time.Now().Add(-1 * time.Minute)) {
		return dto.BookingResponse{}, fmt.Errorf("%w: date in past", ErrInvalidRequest)
	}

	dbServiceType, err := mapServiceType(req.BookingDetails.ServiceType)
	if err != nil {
		return dto.BookingResponse{}, fmt.Errorf("%w: %v", ErrInvalidRequest, err)
	}

	dbTimeSlot, err := mapTimeSlot(req.ScheduleInfo.TimeSlot)
	if err != nil {
		return dto.BookingResponse{}, fmt.Errorf("%w: %v", ErrInvalidRequest, err)
	}

	freq := strings.TrimSpace(req.BookingDetails.Frequency) // в БД: Разовая/Еженедельная/Ежемесячная
	if freq == "" {
		return dto.BookingResponse{}, fmt.Errorf("%w: bad frequency", ErrInvalidRequest)
	}

	// MVP: прайс пока заглушка
	totalPrice := 0.0

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return dto.BookingResponse{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	serviceID, err := s.repo.InsertService(ctx, tx,
		req.BookingDetails.Area,
		dbServiceType,
		freq,
		req.BookingDetails.NoMop,
		req.BookingDetails.NoVacuum,
		req.BookingDetails.HasPet,
	)
	if err != nil {
		return dto.BookingResponse{}, err
	}

	bookingID, createdAt, status, err := s.repo.InsertBooking(ctx, tx,
		serviceID,
		userID,
		scheduled,
		dbTimeSlot,
		req.ScheduleInfo.Address,
		req.ScheduleInfo.Comment,
		totalPrice,
	)
	if err != nil {
		return dto.BookingResponse{}, err
	}

	for _, a := range req.BookingDetails.AdditionalServices {
		if err := s.repo.InsertAdditionalService(ctx, tx, bookingID, a.ID, a.Quantity, 0, 0); err != nil {
			return dto.BookingResponse{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.BookingResponse{}, err
	}

	return dto.BookingResponse{
		ID:         bookingID,
		Status:     status,
		CreatedAt:  createdAt.UTC().Format(time.RFC3339),
		TotalPrice: totalPrice,
	}, nil
}

func mapServiceType(api string) (string, error) {
	switch strings.TrimSpace(api) {
	case "Поддерживающая уборка":
		return "Поддерживающая", nil
	case "Генеральная уборка":
		return "Генеральная", nil
	case "После ремонта":
		return "После ремонта", nil
	default:
		return "", fmt.Errorf("unknown serviceType: %q", api)
	}
}

func mapTimeSlot(api string) (string, error) {
	switch strings.TrimSpace(api) {
	case "09:00 - 12:00":
		return "9-12", nil
	case "12:00 - 15:00":
		return "12-15", nil
	case "15:00 - 18:00":
		return "15-18", nil
	case "18:00 - 21:00":
		return "18-21", nil
	default:
		return "", fmt.Errorf("unknown timeSlot: %q", api)
	}
}
