package model

import "time"

type Vehicle struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement:false"`
	Plate     string    `gorm:"uniqueIndex;size:32;not null"`
	Type      int32     `gorm:"default:1"`
	Status    int32     `gorm:"default:1"`
	Latitude  float64   `gorm:"default:0"`
	Longitude float64   `gorm:"default:0"`
	FuelLevel float64   `gorm:"default:100"`
	MineID    uint64    `gorm:"index"`
	DriverID  uint64    `gorm:"default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (Vehicle) TableName() string { return "vehicles" }
