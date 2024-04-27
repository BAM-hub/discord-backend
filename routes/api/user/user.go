package user

import (
	"database/sql"
	"discord-backend/utils/GetEnv"
	StatusHandler "discord-backend/utils/sql"
	"net/http"

	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID int64 `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID int64 `json:"id"`
	Email string `json:"email"`	
}


var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}
func Encode(b []byte) string {
 return base64.StdEncoding.EncodeToString(b)
}

func Encrypt(text, MySecret string) (string, error) {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
	 return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return Encode(cipherText), nil
}

func CreateUser(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		return 
	}

	if len(newUser.Email) == 0  {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "email is required",
		})
	}
	
	if len(newUser.Password) == 0  {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "password is required",
		})
	}

	secret := GetEnv.GoDotEnvVariable("SECRET")

	encText, err := Encrypt(newUser.Password, secret)
    if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "encrtypt err"})
		return
    }
	db := c.MustGet("db").(*sql.DB)

	stmt, err := db.Prepare("INSERT INTO user(email, password) VALUES(?,?)")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "prepare internal server error"})
		return
	}

	row, execErr := stmt.Exec(newUser.Email, encText)

	if execErr != nil {
		StatusHandler.HandleStatus(err, c)
	}

	defer stmt.Close()

	userId, err := row.LastInsertId()
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "user created successfully", "user":
		  User{ID: userId, Email: newUser.Email, Password: encText}})
}

func GetUsers(c *gin.Context) {
	var users []UserResponse
	db := c.MustGet("db").(*sql.DB)
	rows, err := db.Query("select email, id from user")

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
	} 

	defer rows.Close()

	for rows.Next() {
		var user UserResponse
		if err := rows.Scan(&user.Email, &user.ID); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
		}
		users = append(users, user)
	}

	if rows.Err(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
	}
	c.IndentedJSON(http.StatusOK, users)
}