package controllers

import (
	"net/http"
	"os"
	// "log"
	"time"
	"github.com/Akhilstaar/me-my_encryption/db"
	"github.com/Akhilstaar/me-my_encryption/models"
	"github.com/Akhilstaar/me-my_encryption/utils"
	"github.com/gin-gonic/gin"
	// "github.com/gin-gonic/gin/binding"
)

var Db db.PuppyDb

func AdminLogin(c *gin.Context) {
	info := new(models.AdminLogin)
	if err := c.BindJSON(info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}

	if info.Id != os.Getenv("AdminId"){
		c.JSON(http.StatusForbidden, gin.H{"error" : "This action will be reported."})
		return
	}

	if info.Pass != os.Getenv("AdminPass"){
		c.JSON(http.StatusForbidden, gin.H{"error" : "Invalid Password."})
		return
	}

	token, err := generateJWTToken(info.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
		return
	}
	expirationTime := time.Now().Add(time.Hour * 24)
	cookie := &http.Cookie{
		Name:     "Authorization",
		Value:    token,
		Expires:  expirationTime,
		Path:     "/",
		Domain:   os.Getenv("domain"),
		HttpOnly: true,
		Secure:   false, // Set this to true if you're using HTTPS, false for HTTP
		SameSite: http.SameSiteStrictMode,
	}
	
	http.SetCookie(c.Writer, cookie)
	c.JSON(http.StatusOK, gin.H{"message": "Admin logged in successfully !!"})
}

func AddNewUser(c *gin.Context) {
	// TODO: Modify this function to handle multiple concatenated json inputs


	// TODO: Implement admin authentication logic
	// Authenticate the admin here

	// Validate the input format
	info := new(models.AddNewUser)
	if err := c.BindJSON(info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}

	// Create user
	for _,user := range info.TypeUserNew {

		newUser := models.User{
			Id:      user.Id,
			Name:    user.Name,
			Email:   user.Email,
			Gender:  user.Gender,
			Pass:    user.PassHash,
			PrivK:   "",
			PubK:    "",
			AuthC:   utils.RandStringRunes(15),
			Data:    "",
			Submit:  false,
			Matches: "",
			Dirty:   false,
		}
	
		// Insert the user into the database
	
		if err := Db.Create(&newUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully."})
}

func DeleteUser(c *gin.Context) {
	// TODO: Implement admin authentication logic
	// Authenticate the admin here

	// Validate the input format
	info := new(models.TypeUserNew)
	if err := c.BindJSON(info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}

	newUser := models.User{
		Id:     info.Id,
		Name:   info.Name,
		Email:  info.Email,
		Gender: info.Gender,
	}

	if err := Db.Unscoped().Delete(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User Deleted successfully."})
}

func DeleteAllUsers(c *gin.Context) {
	// TODO: Implement admin authentication logic
	// Authenticate the admin here

	newUser := models.User{}
	if err := Db.Unscoped().Where("1 = 1").Delete(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "All Users Deleted successfully."})
}
