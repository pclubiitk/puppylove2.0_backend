package main

import (

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"

	"github.com/Akhilstaar/me-my_encryption/db"
	"github.com/Akhilstaar/me-my_encryption/router"
	"github.com/Akhilstaar/me-my_encryption/utils"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
	var CfgAdminPass = os.Getenv("CfgAdminPass")

	Db := db.InitDB()

	utils.Randinit()
	store := cookie.NewStore([]byte(CfgAdminPass))
	r := gin.Default()
	r.Use(sessions.Sessions("adminsession", store))
	router.PuppyRoute(r, *Db)

	r.Run(":8080")
	// if err := r.Run(config.CfgAddr); err != nil {
	// 	fmt.Println("[Error] " + err.Error())
	// }
}