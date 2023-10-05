package router

import (
	"net/http"
	"github.com/pclubiitk/puppylove2.0_backend/controllers"
	"github.com/pclubiitk/puppylove2.0_backend/db"
	"github.com/gin-gonic/gin"
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
		users.GET("/fetchall", controllers.FetchHearts)
		// users.GET("/fetchreturnhearts", controllers.FetchReturnHearts)
		users.POST("/sendheart", controllers.SendHeart)
		users.POST("/claimheart", controllers.HeartClaim)
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
