package dbConnect

import (
	"database/sql"
	"discord-backend/utils/GetEnv"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var db *sql.DB

func Connect() (*sql.DB, error) {
	cfg := fmt.Sprintf("%s:%s@%s(%s)/%s", 
	GetEnv.GoDotEnvVariable("DB_USER"),
	GetEnv.GoDotEnvVariable("DB_PASS"),
	"tcp",
	"127.0.0.1:3306",
	"discord",
)
db, err := sql.Open("mysql",cfg)
if err != nil {
	return nil, err
}

println("Connected to database")
return db, nil
}