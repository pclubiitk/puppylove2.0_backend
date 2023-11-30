package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pclubiitk/puppylove2.0_backend/controllers"
	"github.com/pclubiitk/puppylove2.0_backend/db"
)

func PuppyRoute(r *gin.Engine, db db.PuppyDb) {

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, "Hello from the other side!")
	})

	// assigning here cuz I'm only inporting controllers here, & considering their size better import them here.
	controllers.Db = db

	// User administration
	users := r.Group("/users")
	{
		users.POST("/mail/:id", controllers.UserMail)
		users.POST("/login/first", controllers.UserFirstLogin)
		users.Use(controllers.AuthenticateUser())
		users.GET("/fetchPublicKeys", controllers.FetchPublicKeys)
		users.GET("/fetchReturnHearts", controllers.FetchReturnHearts)
		users.POST("/sendheartVirtual", controllers.SendHeartVirtual)
		users.GET("/fetchall", controllers.FetchHearts)
		users.POST("/sentHeartDecoded", controllers.SentHeartDecoded)
		// users.GET("/fetchreturnhearts", controllers.FetchReturnHearts)
		users.POST("/sendheart", controllers.SendHeart)
		users.POST("/claimheart", controllers.HeartClaim)
		users.POST("/publish", controllers.Publish)
	}
	late := r.Group("/special")
	{
		late.Use(controllers.AuthenticateUser())
		late.Use(controllers.AuthenticateUserHeartclaim())
		late.POST("/returnclaimedheartlate", controllers.ReturnClaimedHeartLate)
	}

	// Session administration
	session := r.Group("/session")
	{
		session.POST("/admin/login", controllers.AdminLogin)
		session.POST("/login", controllers.UserLogin)
		session.GET("/logout", controllers.UserLogout)
	}

	// admin
	admin := r.Group("/admin")
	{
		admin.Use(controllers.AuthenticateAdmin())
		admin.GET("/user/deleteallusers", controllers.DeleteAllUsers)
		admin.POST("/user/new", controllers.AddNewUser)
		admin.POST("/user/delete", controllers.DeleteUser)
	}

}
