package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pclubiitk/puppylove2.0_backend/mail"
	"github.com/pclubiitk/puppylove2.0_backend/models"
	"github.com/pclubiitk/puppylove2.0_backend/utils"
	"gorm.io/gorm"
)

func FetchPublicKeys(c *gin.Context) {
	var userModel models.User
	var publicKeys []models.UserPublicKey
	fetchPublicKey := Db.Model(&userModel).Select("id", "pubk").Find(&publicKeys)
	if fetchPublicKey.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Some Error Occured"})
		return
	}
	c.JSON(http.StatusOK, publicKeys)
}
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
func SentHeartDecoded(c *gin.Context) {
	info := new(models.SentHeartsDecoded)
	if err := c.BindJSON(info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}

	var heart models.SendHeart
	var hearts []models.SendHeart

	fetchHeart := Db.Model(&heart).Select("sha", "enc", "gender_of_sender").Find(&hearts)

	if fetchHeart.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No heart to fetch"})
	}

	matchCount := struct {
		male   int
		female int
	}{
		0,
		0,
	}

	for index, heart := range info.DecodedHearts {
		enc := heart.Enc
		gender := heart.GenderOfSender
		if gender != hearts[index].GenderOfSender {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Data not in sequence"})
			return
		}
		if enc == hearts[index].ENC {
			if gender == "M" {
				matchCount.male += 1
			} else {
				matchCount.female += 1
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"male": matchCount.male, "female": matchCount.female})
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
	if u.Dirty {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "User already registered"})
		return
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
