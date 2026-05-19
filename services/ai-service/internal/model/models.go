package model

import "time"

// Query types for AI analysis - maps to existing DB tables

type DispatchTask struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement:false"`
	VehicleID   uint64    `gorm:"not null;index"`
	LoadPointID uint64    `gorm:"not null"`
	DumpPointID uint64    `gorm:"not null"`
	Material    string    `gorm:"size:64"`
	LoadLat     float64   `gorm:"default:0"`
	LoadLon     float64   `gorm:"default:0"`
	DumpLat     float64   `gorm:"default:0"`
	DumpLon     float64   `gorm:"default:0"`
	Status      string    `gorm:"size:32;default:pending"`
	Algorithm   string    `gorm:"size:32;default:fifo"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (DispatchTask) TableName() string { return "dispatch_tasks" }

type Vehicle struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement:false"`
	Plate     string    `gorm:"size:32;uniqueIndex"`
	Type      int32     `gorm:"default:1"`
	Status    int32     `gorm:"default:1"`
	Latitude  float64   `gorm:"default:0"`
	Longitude float64   `gorm:"default:0"`
	FuelLevel float64   `gorm:"default:100"`
	MineID    uint64    `gorm:"default:0"`
	DriverID  uint64    `gorm:"default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (Vehicle) TableName() string { return "vehicles" }

type LoadingPoint struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement:false"`
	Name      string    `gorm:"size:128;not null"`
	Type      string    `gorm:"size:32;not null;default:loading"`
	Latitude  float64   `gorm:"default:0"`
	Longitude float64   `gorm:"default:0"`
	Material  string    `gorm:"size:64"`
	Status    int32     `gorm:"default:1"`
	MineID    uint64    `gorm:"default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (LoadingPoint) TableName() string { return "loading_points" }

type RoadEdge struct {
	ID          uint64  `gorm:"primaryKey;autoIncrement:false"`
	FromNodeID  uint64  `gorm:"not null;index"`
	ToNodeID    uint64  `gorm:"not null;index"`
	DistanceM   float64 `gorm:"not null"`
	MaxSpeedKMH int32   `gorm:"default:30"`
	IsOneway    bool    `gorm:"default:false"`
	MineID      uint64  `gorm:"default:1"`
}

func (RoadEdge) TableName() string { return "road_edges" }

type RoadNode struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement:false"`
	Name      string    `gorm:"size:128;not null"`
	Latitude  float64   `gorm:"not null"`
	Longitude float64   `gorm:"not null"`
	MineID    uint64    `gorm:"default:1"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (RoadNode) TableName() string { return "road_nodes" }
