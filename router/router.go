package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pclubiitk/puppylove2.0_backend/controllers"
	"github.com/pclubiitk/puppylove2.0_backend/db"
)

func PuppyRoute(r *gin.Engine, db db.PuppyDb) {

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Hello from the other side!")
	})

	// assigning here cuz I'm only importing controllers here, & considering their size better import them here.
	controllers.Db = db

	// Captcha
	captcha := r.Group("/captcha/user")
	{
		captcha.Use(controllers.Captchacheck())
		captcha.POST("/mail/:id", controllers.UserMail)
		captcha.POST("/login", controllers.UserLogin)
	}
	// User administration
	users := r.Group("/users")
	{
		// users.POST("/mail/:id", controllers.UserMail)
		users.POST("/retrive", controllers.RetrivePass)
		users.POST("/login/first", controllers.UserFirstLogin)
		users.Use(controllers.AuthenticateUser())
		users.POST("/addRecovery", controllers.AddRecovery)
		users.GET("/data", controllers.GetUserData)
		users.GET("/activeusers", controllers.GetActiveUsers)
		users.GET("/fetchPublicKeys", controllers.FetchPublicKeys)
		users.GET("/fetchReturnHearts", controllers.FetchReturnHearts)

		// profile info
		users.POST("/about", controllers.UpdateAbout)
		users.POST("/intrests", controllers.UpdateIntrest)

		// random search option
		users.GET("/random", controllers.SuggestRandom)

		//Api for verifying hearts that are fetched from Return Table
		users.POST("/verifyreturnhearts", controllers.VerifyReturnHeart)
		users.GET("/fetchall", controllers.FetchHearts)
		users.POST("/sentHeartDecoded", controllers.SentHeartDecoded)
		users.POST("/claimheart", controllers.HeartClaim)

		//API for Last day login
		users.POST("/publish", controllers.Publish)
		users.GET("/mymatches", controllers.MatchesHandler)

		//Send Heart Routes Allowed Only if Admin Permits
		users.Use(controllers.AdminPermit())
		users.POST("/sendheartVirtual", controllers.SendHeartVirtual)
		users.POST("/sendheart", controllers.SendHeart)
	}
	late := r.Group("/special")
	{
		late.Use(controllers.AuthenticateUser())
		users.Use(controllers.AdminPermit())
		// late.Use(controllers.AuthenticateUserHeartclaim())
		late.POST("/returnclaimedheartlate", controllers.ReturnClaimedHeartLate)
	}

	// Session administration
	session := r.Group("/session")
	{
		session.POST("/admin/login", controllers.AdminLogin)
		// session.POST("/login", controllers.UserLogin)
		session.GET("/logout", controllers.UserLogout)
	}

	// admin
	admin := r.Group("/admin")
	{
		admin.Use(controllers.AuthenticateAdmin())
		admin.GET("/user/deleteallusers", controllers.DeleteAllUsers)
		admin.POST("/user/new", controllers.AddNewUser)
		admin.POST("/user/delete", controllers.DeleteUser)
		admin.GET("/publish", controllers.PublishResults)
		admin.GET("/TogglePermit", controllers.TogglePermit)
	}
	r.GET("/stats", controllers.GetStats)

}
