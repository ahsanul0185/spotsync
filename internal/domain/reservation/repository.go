package reservation

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrZoneFull            = errors.New("parking zone is full")
	ErrReservationNotFound = errors.New("reservation not found")
	ErrNotOwner            = errors.New("you can only cancel your own reservations")
)

type Repository interface {
	CreateReservationWithLock(zoneID uint, reservation *Reservation) error
	GetByUserID(userID uint) ([]Reservation, error)
	GetByID(id uint) (*Reservation, error)
	GetAll() ([]Reservation, error)
	CancelReservation(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// CreateReservationWithLock uses GORM transaction with row-level locking (FOR UPDATE)
// to prevent race conditions when multiple users try to book the last available spot.
func (r repository) CreateReservationWithLock(zoneID uint, reservation *Reservation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Lock the parking zone row to prevent concurrent reads/writes
		var zone struct {
			ID            uint
			TotalCapacity int
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Table("parking_zones").
			Select("id, total_capacity").
			Where("id = ?", zoneID).
			Scan(&zone).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("parking zone not found")
			}
			return err
		}

		// 2. Count current active reservations for this zone
		var activeCount int64
		if err := tx.Model(&Reservation{}).
			Where("zone_id = ? AND status = ?", zoneID, "active").
			Count(&activeCount).Error; err != nil {
			return err
		}

		// 3. Check if zone is full
		if int(activeCount) >= zone.TotalCapacity {
			return ErrZoneFull
		}

		// 4. Create the reservation
		if err := tx.Create(reservation).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r repository) GetByUserID(userID uint) ([]Reservation, error) {
	var reservations []Reservation
	if err := r.db.Where("user_id = ?", userID).Preload("Zone").Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r repository) GetByID(id uint) (*Reservation, error) {
	var reservation Reservation
	result := r.db.First(&reservation, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrReservationNotFound
		}
		return nil, result.Error
	}
	return &reservation, nil
}

func (r repository) GetAll() ([]Reservation, error) {
	var reservations []Reservation
	if err := r.db.Preload("User").Preload("Zone").Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r repository) CancelReservation(id uint) error {
	return r.db.Model(&Reservation{}).Where("id = ?", id).Update("status", "cancelled").Error
}
