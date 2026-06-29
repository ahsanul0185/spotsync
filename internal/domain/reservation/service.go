package reservation

import (
	"spotsync/internal/domain/parkingzone"
	"spotsync/internal/domain/reservation/dto"
	"spotsync/internal/domain/user"
)

var (
	ErrZoneNotFound = fmt.Errorf("parking zone not found")
)

type Service interface {
	CreateReservation(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error)
	GetMyReservations(userID uint) ([]dto.ReservationResponse, error)
	CancelReservation(userID uint, reservationID uint) error
	GetAllReservations() ([]dto.ReservationResponse, error)
}

type service struct {
	repo        Repository
	zoneRepo    parkingzone.Repository
	userRepo    user.Repository
}

func NewService(repo Repository, zoneRepo parkingzone.Repository, userRepo user.Repository) Service {
	return &service{
		repo:     repo,
		zoneRepo: zoneRepo,
		userRepo: userRepo,
	}
}

func (s *service) CreateReservation(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	// Verify zone exists
	zone, err := s.zoneRepo.GetByID(req.ZoneID)
	if err != nil {
		if err == parkingzone.ErrZoneNotFound {
			return nil, ErrZoneNotFound
		}
		return nil, err
	}

	reservation := Reservation{
		UserID:       userID,
		ZoneID:       req.ZoneID,
		LicensePlate: req.LicensePlate,
		Status:       "active",
	}

	if err := s.repo.CreateReservationWithLock(req.ZoneID, &reservation); err != nil {
		if err == ErrZoneFull {
			return nil, ErrZoneFull
		}
		return nil, err
	}

	// Populate zone info for response
	reservation.Zone = *zone

	return s.toResponse(&reservation), nil
}

func (s *service) GetMyReservations(userID uint) ([]dto.ReservationResponse, error) {
	reservations, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ReservationResponse, len(reservations))
	for i, r := range reservations {
		responses[i] = *s.toResponse(&r)
	}
	return responses, nil
}

func (s *service) CancelReservation(userID uint, reservationID uint) error {
	reservation, err := s.repo.GetByID(reservationID)
	if err != nil {
		return err
	}

	if reservation.UserID != userID {
		return ErrNotOwner
	}

	if reservation.Status == "cancelled" {
		return nil // already cancelled
	}

	return s.repo.CancelReservation(reservationID)
}

func (s *service) GetAllReservations() ([]dto.ReservationResponse, error) {
	reservations, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ReservationResponse, len(reservations))
	for i, r := range reservations {
		responses[i] = *s.toResponse(&r)
	}
	return responses, nil
}

func (s *service) toResponse(r *Reservation) *dto.ReservationResponse {
	resp := &dto.ReservationResponse{
		ID:           r.ID,
		UserID:       r.UserID,
		ZoneID:       r.ZoneID,
		LicensePlate: r.LicensePlate,
		Status:       r.Status,
		CreatedAt:    r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    r.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	// Populate zone info if available
	if r.Zone.ID != 0 {
		resp.Zone = &dto.ZoneSummary{
			ID:   r.Zone.ID,
			Name: r.Zone.Name,
			Type: r.Zone.Type,
		}
	}

	// Populate user info if available
	if r.User.ID != 0 {
		resp.User = &dto.UserSummary{
			ID:    r.User.ID,
			Name:  r.User.Name,
			Email: r.User.Email,
		}
	}

	return resp
}
