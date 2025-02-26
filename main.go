package main

import (
	"os"

	// "github.com/gin-contrib/cors"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pclubiitk/puppylove2.0_backend/db"
	"github.com/pclubiitk/puppylove2.0_backend/router"
	"github.com/pclubiitk/puppylove2.0_backend/utils"
	"github.com/pclubiitk/puppylove2.0_backend/redisclient"
	"github.com/JGLTechnologies/gin-rate-limit"
	"time"
	"strconv"
)
func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.String(429, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
}
func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
	var CfgAdminPass = os.Getenv("CFG_ADMIN_PASS")
	Db := db.InitDB()
	redisclient.InitRedis()
	utils.Randinit()
	store := cookie.NewStore([]byte(CfgAdminPass))
	rateLimitStr := os.Getenv("RATE_LIMIT")
	rateLimit, err := strconv.ParseUint(rateLimitStr, 10, 32)
	if err != nil {
		panic("Invalid RATE_LIMIT value in .env file")
	}
	rateLimitStore := ratelimit.RedisStore(&ratelimit.RedisOptions{
			RedisClient: redisclient.RedisClient, 
			Rate:        time.Second,       
			Limit:       uint(rateLimit),
		})
	rateLimitMiddleware := ratelimit.RateLimiter(rateLimitStore, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})
	r := gin.Default()
	r.Use(rateLimitMiddleware)
	// Local Testing Frontend Running on localhost:3000
	// r.Use(cors.New(cors.Config{AllowCredentials: true, AllowOrigins: []string{"http://localhost:3000"}, AllowHeaders: []string{"content-type"}}))
	// Allow all origins
	r.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowHeaders:     []string{"content-type", "g-recaptcha-response"},
	}))
	r.Use(sessions.Sessions("adminsession", store))
	router.PuppyRoute(r, *Db)

	r.Run(":8080")

	// if err := r.Run(config.CfgAddr); err != nil {
	// 	fmt.Println("[Error] " + err.Error())
	// }
}
