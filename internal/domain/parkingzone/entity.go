package parkingzone

import (
	"gorm.io/gorm"
)

type ParkingZone struct {
	gorm.Model
	Name          string  `json:"name" gorm:"type:varchar(255);not null"`
	Type          string  `json:"type" gorm:"type:varchar(50);not null"`
	TotalCapacity int     `json:"total_capacity" gorm:"not null"`
	PricePerHour  float64 `json:"price_per_hour" gorm:"not null"`
}
