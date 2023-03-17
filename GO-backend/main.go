package main

import (
	"fmt"
	"net/http"
	AuthController "se/jwt-api/controller/auth"
	"se/jwt-api/controller/middleware"
	CarController "se/jwt-api/controller/car"
	UserController "se/jwt-api/controller/user"
	BookingController "se/jwt-api/controller/booking"
	"se/jwt-api/orm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin" // เป็ น framework ช่วยในการสร้าง api และประสิทธิภาพ
	"github.com/joho/godotenv"
	"gorm.io/gorm" // framwork ต่อกับ database ภาษา GO
)

type Register struct { // ออกแบบข้อมูล
	Username string
	Password string
	Fullname string
}
type User struct { // สร้าง ตารางใน database ชื่อ User
	gorm.Model
	Username string
	Password string
	Fullname string
}
type Car struct {
	gorm.Model
	carname string
	detail  string
	image   string
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	orm.InitDB()
	r := gin.Default()
	r.Use(cors.Default()) // เพื่อให้สามารถเรียก api เราได้
	r.POST("/register", AuthController.Register)
	r.POST("/login", AuthController.Login)
	authorized := r.Group("/users", middleware.JWTAuth())
	authorized.GET("/readall", UserController.ReadAll)

	// api each others
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		var user orm.User
		orm.Db.First(&user, id)
		c.JSON(200, user)
	})

	// soft delete user id
	r.DELETE("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		var user orm.User
		orm.Db.First(&user, id)
		orm.Db.Delete(&user)
		c.JSON(200, user)
	})

	// update user id
	r.PUT("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		var user orm.User
		var updateUser orm.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		orm.Db.First(&updateUser, id)
		updateUser.Username = user.Username
		updateUser.Fullname = user.Fullname
		orm.Db.Save(updateUser)
		c.JSON(200, updateUser)
	})


	// api set of car
	r.POST("/carregister", CarController.RegisterCar)
	r.GET("/carall", CarController.CarAll)
	// api each others
	r.GET("/cars/:id", func(c *gin.Context) {
		id := c.Param("id")
		var car orm.Car
		orm.Db.First(&car, id)
		c.JSON(200, car)
	})
	//Update
	r.PUT("/cars/:id", func(c *gin.Context) {
		id := c.Param("id")
		var car orm.Car
		var updateCar orm.Car
		if err := c.ShouldBindJSON(&car); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		orm.Db.First(&updateCar, id)
		updateCar.Carname = car.Carname
		updateCar.Detail = car.Detail
		updateCar.Image = car.Image
		orm.Db.Save(updateCar)
		c.JSON(200, updateCar)
	})
	// soft delete user id
	r.DELETE("/cars/:id", func(c *gin.Context) {
		id := c.Param("id")
		var car orm.Car
		orm.Db.First(&car, id)
		orm.Db.Delete(&car)
		c.JSON(200, car)
	})
	//Api Check Booking
	r.POST("/bookingcar",BookingController.BookingCar)

	r.Run("localhost:8080")

}
