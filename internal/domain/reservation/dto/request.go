package dto

type CreateReservationRequest struct {
	ZoneID       uint   `json:"zone_id" validate:"required,gt=0"`
	LicensePlate string `json:"license_plate" validate:"required,max=15"`
}

type ReservationResponse struct {
	ID           uint             `json:"id"`
	UserID       uint             `json:"user_id"`
	ZoneID       uint             `json:"zone_id"`
	LicensePlate string           `json:"license_plate"`
	Status       string           `json:"status"`
	Zone         *ZoneSummary     `json:"zone,omitempty"`
	User         *UserSummary     `json:"user,omitempty"`
	CreatedAt    string           `json:"created_at"`
	UpdatedAt    string           `json:"updated_at"`
}

type ZoneSummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type UserSummary struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
