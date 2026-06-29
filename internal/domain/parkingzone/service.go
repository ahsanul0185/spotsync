package parkingzone

import (
	"spotsync/internal/domain/parkingzone/dto"
)

type Service interface {
	Create(req dto.CreateZoneRequest) (*dto.ZoneResponse, error)
	GetAll() ([]dto.ZoneResponse, error)
	GetByID(id uint) (*dto.ZoneResponse, error)
	Update(id uint, req dto.UpdateZoneRequest) (*dto.ZoneResponse, error)
	Delete(id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(req dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
	zone := ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}
	if err := s.repo.Create(&zone); err != nil {
		return nil, err
	}
	return s.toResponse(&zone), nil
}

func (s *service) GetAll() ([]dto.ZoneResponse, error) {
	zones, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ZoneResponse, len(zones))
	for i, zone := range zones {
		responses[i] = *s.toResponse(&zone)
	}
	return responses, nil
}

func (s *service) GetByID(id uint) (*dto.ZoneResponse, error) {
	zone, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(zone), nil
}

func (s *service) Update(id uint, req dto.UpdateZoneRequest) (*dto.ZoneResponse, error) {
	zone, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		zone.Name = req.Name
	}
	if req.Type != "" {
		zone.Type = req.Type
	}
	if req.TotalCapacity > 0 {
		zone.TotalCapacity = req.TotalCapacity
	}
	if req.PricePerHour > 0 {
		zone.PricePerHour = req.PricePerHour
	}

	if err := s.repo.Update(zone); err != nil {
		return nil, err
	}
	return s.toResponse(zone), nil
}

func (s *service) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *service) toResponse(zone *ParkingZone) *dto.ZoneResponse {
	activeCount, _ := s.repo.CountActiveReservations(zone.ID)
	available := zone.TotalCapacity - int(activeCount)
	if available < 0 {
		available = 0
	}

	return &dto.ZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: available,
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:      zone.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
