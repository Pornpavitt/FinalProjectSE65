package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"se/jwt-api/orm"
)

var hmacSampleSecret []byte

func ReadAll(c *gin.Context) {
	var users []orm.User
	orm.Db.Find(&users)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Read Sucessful",
		"users": users})
}
