package dto

type AdditionalService struct {
	ID       string `json:"id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,min=1"`
}

type BookingDetails struct {
	ServiceType        string              `json:"serviceType" binding:"required"`
	Area               int                 `json:"area" binding:"required,min=10"`
	Frequency          string              `json:"frequency" binding:"required"`
	NoMop              bool                `json:"noMop"`
	NoVacuum           bool                `json:"noVacuum"`
	HasPet             bool                `json:"hasPet"`
	AdditionalServices []AdditionalService `json:"additionalServices"`
}

type ScheduleInfo struct {
	Date     string `json:"date" binding:"required"`
	TimeSlot string `json:"timeSlot" binding:"required"`
	Address  string `json:"address" binding:"required"`
	Comment  string `json:"comment"`
}

type AuthBookingRequest struct {
	BookingDetails BookingDetails `json:"bookingDetails" binding:"required"`
	ScheduleInfo   ScheduleInfo   `json:"scheduleInfo" binding:"required"`
}

type BookingResponse struct {
	ID         string  `json:"id"`
	Status     string  `json:"status"`
	CreatedAt  string  `json:"createdAt"`
	TotalPrice float64 `json:"totalPrice"`
}
