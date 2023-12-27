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
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	user := os.Getenv("POSTGRES_USER")

	// if os.Getenv(keyDocker)=""{
	// 	host = "localhost"
	// }

	loginstring := fmt.Sprintf("host=%s user=%s  password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata", host, user, password, dbName, port)

	db, err := gorm.Open(postgres.Open(loginstring), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	db.AutoMigrate(&models.User{}, &models.SendHeart{}, &models.HeartClaims{}, &models.ReturnHearts{}, &models.MatchTable{})
	fmt.Println("Connected to the database!")
	// sqlDB.Close()
	return &PuppyDb{db}
}
