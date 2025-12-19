package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type BookingRepo struct{}

func NewBookingRepo() *BookingRepo { return &BookingRepo{} }

func (r *BookingRepo) InsertService(ctx context.Context, tx pgx.Tx, area int, serviceType string, frequency string, noMop, noVacuum, hasPet bool) (int, error) {
	var id int
	err := tx.QueryRow(ctx, `
		insert into services (area, service_type, frequency, no_mop, no_vacuum, has_pet)
		values ($1, $2::service_type, $3::frequency_type, $4, $5, $6)
		returning id
	`, area, serviceType, frequency, noMop, noVacuum, hasPet).Scan(&id)
	return id, err
}

func (r *BookingRepo) InsertBooking(ctx context.Context, tx pgx.Tx, serviceID int, userID string, scheduled time.Time, timeSlot string, address string, comment string, totalPrice float64) (bookingID string, createdAt time.Time, status string, err error) {
	err = tx.QueryRow(ctx, `
		insert into bookings (service_id, user_id, scheduled_date, time_slot, address, comment, total_price)
		values ($1, $2::uuid, $3, $4::time_slot_type, $5, $6, $7)
		returning id, created_at, status::text
	`, serviceID, userID, scheduled, timeSlot, address, comment, totalPrice).Scan(&bookingID, &createdAt, &status)

	return
}

func (r *BookingRepo) InsertAdditionalService(ctx context.Context, tx pgx.Tx, bookingID string, serviceType string, qty int, unitPrice float64, totalPrice float64) error {
	_, err := tx.Exec(ctx, `
		insert into booking_additional_services (booking_id, service_type, quantity, unit_price, total_price)
		values ($1::uuid, $2::additional_service_type, $3, $4, $5)
	`, bookingID, serviceType, qty, unitPrice, totalPrice)
	return err
}
