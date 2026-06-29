package parkingzone

import (
	"errors"

	"gorm.io/gorm"
)

var ErrZoneNotFound = errors.New("parking zone not found")

type Repository interface {
	Create(zone *ParkingZone) error
	GetAll() ([]ParkingZone, error)
	GetByID(id uint) (*ParkingZone, error)
	Update(zone *ParkingZone) error
	Delete(id uint) error
	CountActiveReservations(zoneID uint) (int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r repository) Create(zone *ParkingZone) error {
	return r.db.Create(zone).Error
}

func (r repository) GetAll() ([]ParkingZone, error) {
	var zones []ParkingZone
	if err := r.db.Find(&zones).Error; err != nil {
		return nil, err
	}
	return zones, nil
}

func (r repository) GetByID(id uint) (*ParkingZone, error) {
	var zone ParkingZone
	result := r.db.First(&zone, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrZoneNotFound
		}
		return nil, result.Error
	}
	return &zone, nil
}

func (r repository) Update(zone *ParkingZone) error {
	return r.db.Save(zone).Error
}

func (r repository) Delete(id uint) error {
	return r.db.Delete(&ParkingZone{}, id).Error
}

func (r repository) CountActiveReservations(zoneID uint) (int64, error) {
	var count int64
	err := r.db.Table("reservations").
		Where("zone_id = ? AND status = ?", zoneID, "active").
		Count(&count).Error
	return count, err
}
