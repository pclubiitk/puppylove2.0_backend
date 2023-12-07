package main

import (
	"os"

	// "github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pclubiitk/puppylove2.0_backend/db"
	"github.com/pclubiitk/puppylove2.0_backend/router"
	"github.com/pclubiitk/puppylove2.0_backend/utils"
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
	// r.Use(cors.Default())
	r.Use(corsMiddleware())
	router.PuppyRoute(r, *Db)

	r.Run(":8080")

	// if err := r.Run(config.CfgAddr); err != nil {
	// 	fmt.Println("[Error] " + err.Error())
	// }
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Set the origin of your frontend app.
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true") // Allow credentials for preflight requests.
			c.AbortWithStatus(204)
			return
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true") // Allow credentials for the main request.
		c.Next()
	}
}
