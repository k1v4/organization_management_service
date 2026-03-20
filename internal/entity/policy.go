package entity

type UpdatePolicy struct {
	MaxBookingDurationMin    *int `json:"max_booking_duration_min"`
	BookingWindowDays        *int `json:"booking_window_days"`
	MaxActiveBookingsPerUser *int `json:"max_active_bookings_per_user"`
}
