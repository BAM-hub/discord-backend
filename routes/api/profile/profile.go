package profile

import (
	"database/sql"
	StatusHandler "discord-backend/utils/sql"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)


type Profile struct {
	Id int64 `json:"id"`
	UserId int64 `json:"userId"`
	Avatar string `json:"avatar"`
	DisplayName string `json:"displayName"`
	UserName string `json:"userName"`
	PhoneNumber string `json:"phoneNumber"`
	Status string `json:"status"`
	CustomStatus *string `json:"customStatus"`
	ClearAfter *time.Time `json:"clearAfter"`
}

func CreateProfile(c *gin.Context) {
	var profile Profile
	err :=	c.BindJSON(&profile)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, 
			gin.H{"message": "request parse error!",
		})
		return
	}

	if Type := reflect.TypeOf(profile.UserId).Kind(); Type ==  reflect.Int  {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "id is required", 
		})	
		return
	}

	if len(profile.DisplayName) == 0 || len(profile.UserName) == 0 { 
		message := fmt.Sprintf("field %s %s is required",
			"displayName",
			"userName",
		)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": message, 
		})	
		return
	}

	db := c.MustGet("db").(*sql.DB)

	stmt, err := db.Prepare("INSERT INTO profile(userName, displayName, avatar, phoneNumber, userId) VALUES(?,?,?,?,?)")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error while prepare stmt",
			"err": err.Error(),
		})
		return
	}

	defer stmt.Close()

	result, err :=  stmt.Exec(profile.UserName, profile.DisplayName, profile.Avatar, profile.PhoneNumber, profile.UserId)
	if err != nil {
		StatusHandler.HandleStatus(err, c)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}	

	c.IndentedJSON(http.StatusOK, gin.H{
		"id": id,
		"userName": profile.UserName,
		"displayName": profile.DisplayName,
		"avatar": profile.Avatar,
		"phoneNumber": profile.PhoneNumber,
		"status": "offline",
	})

}

type GetProfileByIdRequest struct {
	ProfileId int64 `json:"profileId"`
}

func GetProfileById(c *gin.Context) {
	vars := c.Query("profileId")
	id, err := strconv.ParseInt(vars, 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "excpectd profileId as number",
		})
	}
	var profile Profile


	db := c.MustGet("db").(*sql.DB)
	row := db.QueryRow(
		"select  id, avatar, displayName, userName, phoneNumber, status, customStatus, clearAfter, userId from profile where id = ?",
		id,
	)
	

	if err = row.Scan(&profile.Id, &profile.Avatar, &profile.DisplayName, &profile.UserName, &profile.PhoneNumber, &profile.Status, &profile.CustomStatus, &profile.ClearAfter, &profile.UserId); err != nil {				
		println("error", err.Error())
		StatusHandler.HandleStatus(err, c)
		return
	}

	c.JSON(http.StatusOK, profile)
}

type GetProfileByUserIdRequest struct {
	UserId int64 `json:"userId"`
}

func GetProfileByUserId(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "excpectd userId as number",
		})
	}
	var profile Profile

	db := c.MustGet("db").(*sql.DB)
	row := db.QueryRow(
		"select  id, avatar, displayName, userName, phoneNumber, status, customStatus, clearAfter, userId from profile where userId = ?",
		id,
	)
	
	if err =  row.Scan(&profile); err != nil {
		StatusHandler.HandleStatus(err, c)
		return
	}
	c.JSON(http.StatusOK,  profile)

}