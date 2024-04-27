package main

import (
	dbConnect "discord-backend/db"
	"discord-backend/routes/api/profile"
	user "discord-backend/routes/api/user"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
  db, err := dbConnect.Connect()
  if err != nil {
    log.Fatalf("Could not connect to the database: %v", err)
  }
  router := gin.Default()
  router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	})) 
  router.Use(func(c *gin.Context) {
    c.Set("db", db)
    c.Next()
  })
  router.POST("/user/create", user.CreateUser)
  router.GET("/user/list", user.GetUsers)
  router.POST("/profile/create", profile.CreateProfile)
  router.Run("localhost:8080")
}

