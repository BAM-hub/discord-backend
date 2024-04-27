package StatusHandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func HandleStatus(err error, c *gin.Context) {
	errNumber := err.(*mysql.MySQLError).Number
	if errNumber == 1452 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "record refecernced dose not exist!",
		})
	} else if errNumber == 1062 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
	} else if errNumber == 1364 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "required field missing"}) 
	} else {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error result err",
		})
	}
}