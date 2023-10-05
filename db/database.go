package db

import (
	"fmt"
	"os"
	"github.com/pclubiitk/puppylove2.0_backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type PuppyDb struct {
	*gorm.DB
}

func InitDB() *PuppyDb {
	host := os.Getenv("host")
	port := os.Getenv("port")
	password := os.Getenv("password")
	dbName := os.Getenv("dbName")
	user := os.Getenv("user")

	loginstring := fmt.Sprintf("host=%s user=%s  password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata", host, user, password, dbName, port)

	db, err := gorm.Open(postgres.Open(loginstring), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	db.AutoMigrate(&models.User{}, &models.SendHeart{}, &models.HeartClaims{}, &models.ReturnHearts{})
	fmt.Println("Connected to the database!")
	// sqlDB.Close()
	return &PuppyDb{db}
}
