package controllers

import (
	"errors"
	"net/http"

	"github.com/Akhilstaar/me-my_encryption/mail"
	"github.com/Akhilstaar/me-my_encryption/models"
	"github.com/Akhilstaar/me-my_encryption/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FetchHearts(c *gin.Context) {
	var heart models.SendHeart
	var hearts []models.FetchHeartsFirst
	// Fetch only required columns from the database
	fetchheart := Db.Model(&heart).Select("enc", "gender_of_sender").Find(&hearts)

	if fetchheart.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No hearts to fetch."})
		return
	}

	c.JSON(http.StatusOK, hearts)
}

func UserMail(c *gin.Context) {
	id := c.Param("id")
	u := models.MailData{}
	user := models.User{}
	record := Db.Model(&user).Where("id = ?", id).First(&u)
	if record.Error != nil {
		if errors.Is(record.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{"error": "User not found !!"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
			return
		}
	}
	AuthC := utils.RandStringRunes(15)
	Db.Model(&user).Where("id = ?", id).Update("AuthC", AuthC)
	if mail.SendMail(u.Name, u.Email, AuthC) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Auth. code sent successfully !!"})
}

// func FetchReturnHearts(c *gin.Context) {
// 	heartModel := models.ReturnHearts{}
//     var hearts []models.ReturnHearts
// 	fetchheart := Db.Model(&heartModel).Select("enc").Find(hearts)
// 	if fetchheart.Error != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error" : "No hearts to fetch."})
//         return
// 	}

// 	c.JSON(http.StatusOK, hearts)
// }
