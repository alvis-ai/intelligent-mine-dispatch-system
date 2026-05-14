package model

import "time"

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

type RoadEdge struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement:false"`
	FromNodeID uint64    `gorm:"not null;index"`
	ToNodeID   uint64    `gorm:"not null;index"`
	DistanceM  float64   `gorm:"not null"`
	MaxSpeedKMH int32    `gorm:"default:30"`
	IsOneway   bool      `gorm:"default:false"`
	MineID     uint64    `gorm:"default:1"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func (RoadEdge) TableName() string { return "road_edges" }
