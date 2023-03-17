package booking

import (
	"net/http"
	"se/jwt-api/orm"
	"time"

	"github.com/gin-gonic/gin"
)

// สร้าง Structure เพื่อ รองรับ Json
type BookingBody struct {
	ID	   uint
	UserID uint
	CarID  uint
	Start  time.Time
	End    time.Time
}

func BookingCar(c *gin.Context) {
	var json BookingBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	layout := "2006-01-02"
	start, err := time.Parse(layout, json.Start.Format(layout))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	end, err := time.Parse(layout, json.End.Format(layout))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if end.Before(start) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End time is before Start time"})
		return
	}
	// Check if the car is already booked during the requested time period
	var bookings []orm.Booking
	orm.Db.Where("car_id = ? AND ((start <= ? AND ?) OR (end >= ? AND ?))", json.CarID, start, end, start, end).Find(&bookings)

	for _, b := range bookings {
		if b.ID != 0 || b.CarID != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Car is already booked during requested time period"})
			return
		}
	}

	// Query the database using Gorm
	var results []orm.Booking
	orm.Db.Where("car_id = ? AND start BETWEEN ? AND ?", json.CarID, start, end).Find(&results)

	// Check if the booking already exists
	if len(results) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking already exists"})
		return
	}

	// Create the booking
	booking := orm.Booking{UserID: json.UserID, CarID: json.CarID, Start: start, End: end}
	if err := orm.Db.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": booking})
}
