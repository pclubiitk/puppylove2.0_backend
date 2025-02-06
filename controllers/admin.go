package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pclubiitk/puppylove2.0_backend/db"
	"github.com/pclubiitk/puppylove2.0_backend/models"
	"github.com/pclubiitk/puppylove2.0_backend/utils"
)

var Db db.PuppyDb
var permit bool = true

func AdminLogin(c *gin.Context) {
	info := new(models.AdminLogin)
	if err := c.BindJSON(info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}

	if info.Id != os.Getenv("ADMIN_ID") {
		c.JSON(http.StatusForbidden, gin.H{"error": "This action will be reported."})
		return
	}

	if info.Pass != os.Getenv("ADMIN_PASS") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Password."})
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
		Domain:   os.Getenv("DOMAIN"),
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
	for _, user := range info.TypeUserNew {

		emptyMatches, _ := json.Marshal(make(map[string]string))
		newUser := models.User{
			Id:        user.Id,
			Name:      user.Name,
			Email:     user.Email,
			Gender:    user.Gender,
			Pass:      "",
			PubK:      "",
			PrivK:     "",
			AuthC:     utils.RandStringRunes(15),
			Data:      "",
			Submit:    false,
			Matches:   emptyMatches,
			Dirty:     false,
			Publish:   false,
			Code:      "",
			About:     "",
			Interests: "",
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

func PublishResults(c *gin.Context) {
	if !models.PublishMatches {
		var matchdb models.MatchTable
		var matches []models.MatchTable
		records := Db.Model(&matchdb).Find(&matches)
		if records.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Some error occurred while calculating matches"})
			return
		}

		for _, match := range matches {
			roll1 := match.Roll1
			roll2 := match.Roll2
			song12 := match.SONG12 // Song sent by roll2 to roll1
			song21 := match.SONG21 // Song sent by roll1 to roll2

			var userdb models.User
			var userdb1 models.User

			// Fetch user1
			if err := Db.Model(&userdb).Where("id = ?", roll1).First(&userdb).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching user1"})
				return
			}

			// Fetch user2
			if err := Db.Model(&userdb1).Where("id = ?", roll2).First(&userdb1).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching user2"})
				return
			}

			// Only proceed if both users have opted to publish their results
			if userdb.Publish && userdb1.Publish {
				// Unmarshal existing matches
				var matches1, matches2 map[string]string
				if len(userdb.Matches) > 0 {
					_ = json.Unmarshal(userdb.Matches, &matches1)
				} else {
					matches1 = make(map[string]string)
				}

				if len(userdb1.Matches) > 0 {
					_ = json.Unmarshal(userdb1.Matches, &matches2)
				} else {
					matches2 = make(map[string]string)
				}

				// Add new matches
				matches1[roll2] = song21
				matches2[roll1] = song12

				// Marshal back into json.RawMessage
				userdb.Matches, _ = json.Marshal(matches1)
				userdb1.Matches, _ = json.Marshal(matches2)

				// Save the updated user records
				Db.Save(&userdb)
				Db.Save(&userdb1)
			}
		}

		models.PublishMatches = true
		c.JSON(http.StatusOK, gin.H{"msg": "Published Matches"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Matches already published"})
}

func TogglePermit(c *gin.Context) {
	permit = !permit
	c.JSON(http.StatusOK, gin.H{"permitStatus": permit})
}
