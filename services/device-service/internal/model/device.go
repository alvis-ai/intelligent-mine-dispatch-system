package model

import "time"

type Device struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement:false"`
	Name           string    `gorm:"size:128;not null"`
	DeviceType     int32     `gorm:"default:0"`
	Status         int32     `gorm:"default:1"`
	FirmwareVersion string   `gorm:"size:64;default:''"`
	Latitude       float64   `gorm:"default:0"`
	Longitude      float64   `gorm:"default:0"`
	MineID         uint64    `gorm:"index;default:0"`
	VehicleID      uint64    `gorm:"default:0"`
	LastOnlineAt   time.Time `gorm:""`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

func (Device) TableName() string { return "devices" }
