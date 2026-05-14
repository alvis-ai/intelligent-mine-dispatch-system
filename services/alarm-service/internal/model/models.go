package model

import "time"

type Geofence struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement:false"`
	Name        string    `gorm:"size:128;not null"`
	Shape       string    `gorm:"size:16;default:'circle'"`
	CenterLat   float64   `gorm:"default:0"`
	CenterLon   float64   `gorm:"default:0"`
	RadiusM     float64   `gorm:"default:0"`
	PointsJSON  string    `gorm:"column:points_json;type:text"`
	FenceType   string    `gorm:"size:32;default:'restricted'"`
	MinSpeedKMH int32     `gorm:"default:0"`
	MaxSpeedKMH int32     `gorm:"default:0"`
	TimeRange   string    `gorm:"size:32"`
	Enabled     bool      `gorm:"default:true"`
	MineID      uint64    `gorm:"default:1"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (Geofence) TableName() string { return "geofences" }

type AlarmRule struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement:false"`
	Name        string    `gorm:"size:128;not null"`
	RuleType    string    `gorm:"size:32;not null"`
	GeofenceID  uint64    `gorm:"default:0"`
	Severity    string    `gorm:"size:16;default:'warning'"`
	Description string    `gorm:"size:256"`
	Enabled     bool      `gorm:"default:true"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (AlarmRule) TableName() string { return "alarm_rules" }

type AlarmEvent struct {
	ID              uint64    `gorm:"primaryKey;autoIncrement:false"`
	RuleID          uint64    `gorm:"index"`
	VehicleID       uint64    `gorm:"index"`
	VehiclePlate    string    `gorm:"size:64"`
	AlarmType       string    `gorm:"size:32;not null"`
	Severity        string    `gorm:"size:16;default:'warning'"`
	Message         string    `gorm:"size:512"`
	Latitude        float64   `gorm:"default:0"`
	Longitude       float64   `gorm:"default:0"`
	Speed           float64   `gorm:"default:0"`
	Acknowledged    bool      `gorm:"default:false"`
	AcknowledgedBy  string    `gorm:"size:64"`
	AcknowledgedAt  *time.Time
	MineID          uint64    `gorm:"default:1"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
}

func (AlarmEvent) TableName() string { return "alarm_events" }

const (
	SeverityInfo     = "info"
	SeverityWarning  = "warning"
	SeverityCritical = "critical"

	RuleTypeGeofence = "geofence"
	RuleTypeSpeeding = "speeding"
	RuleTypeOffline  = "offline"
	RuleTypeDeviation = "deviation"

	ShapeCircle  = "circle"
	ShapePolygon = "polygon"
)
