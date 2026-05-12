package model

import "time"

type DispatchTask struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement:false"`
	VehicleID   uint64    `gorm:"index"`
	LoadPointID uint64    `gorm:"not null"`
	DumpPointID uint64    `gorm:"not null"`
	Material    string    `gorm:"size:64"`
	LoadLat     float64
	LoadLon     float64
	DumpLat     float64
	DumpLon     float64
	Status      string    `gorm:"size:32;default:'pending'"`
	Algorithm   string    `gorm:"size:32"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (DispatchTask) TableName() string { return "dispatch_tasks" }

const (
	StatusPending   = "pending"
	StatusActive    = "active"
	StatusCompleted = "completed"
	StatusCancelled = "cancelled"
)
