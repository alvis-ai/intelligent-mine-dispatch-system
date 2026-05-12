package model

import "time"

type User struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement:false"`
	Username  string    `gorm:"uniqueIndex;size:64;not null"`
	Password  string    `gorm:"size:256;not null"`
	RealName  string    `gorm:"size:64"`
	Email     string    `gorm:"size:128"`
	Phone     string    `gorm:"size:20"`
	Role      int32     `gorm:"default:0"`
	Status    int32     `gorm:"default:1"`
	MineID    uint64    `gorm:"index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}
