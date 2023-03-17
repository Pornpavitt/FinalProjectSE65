package orm

import (
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	UserID    uint
	CarID     uint
	Start     time.Time
	End       time.Time
}

