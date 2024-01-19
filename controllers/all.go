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
	fetchPublicKey := Db.Model(&userModel).Select("id, pub_k").Find(&publicKeys)
	if fetchPublicKey.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Some Error Occured"})
		return
	}
	c.JSON(http.StatusOK, publicKeys)
}

func FetchReturnHearts(c *gin.Context) {
	var returnedHeart models.ReturnHearts
	var returnedHearts []models.FetchReturnedHearts

	fetchedReturnedHearts := Db.Model(&returnedHeart).Select("enc").Find(&returnedHearts)

	if fetchedReturnedHearts.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No returend hearts to fetch."})
		return
	}

	c.JSON(http.StatusOK, returnedHearts)
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

func ForgotMail(c *gin.Context) {
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

func GetStats(c *gin.Context) {
	var userdb models.User
	var users []models.User

	records := Db.Model(&userdb).Where("dirty = ?", true).Find(&users)
	if records.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Fetching Stats"})
		return
	}

	var matchdb models.MatchTable
	var matches []models.MatchTable

	records = Db.Model(&matchdb).Where("").Find(&matches)
	if records.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Fetching Stats"})
		return
	}
	femaleRegisters := 0
	maleRegisters := 0
	registers := struct {
		y19 int
		y20 int
		y21 int
		y22 int
		y23 int
	}{y19: 0, y20: 0, y21: 0, y22: 0, y23: 0}
	numberOfMatches := int(len(matches) / 2)
	batchwiseMatches := struct {
		y19 int
		y20 int
		y21 int
		y22 int
		y23 int
	}{y19: 0, y20: 0, y21: 0, y22: 0, y23: 0}
	for _, user := range users {
		if user.Gender == "M" {
			maleRegisters++
		} else {
			femaleRegisters++
		}
		startstring := user.Id[0:2]
		switch startstring {
		case "19":
			registers.y19++
		case "20":
			registers.y20++
		case "21":
			registers.y21++
		case "22":
			registers.y22++
		case "23":
			registers.y23++
		}
	}

	for _, matc := range matches {
		startstr := matc.Roll1[0:2]
		switch startstr {
		case "19":
			batchwiseMatches.y19++
		case "20":
			batchwiseMatches.y20++
		case "21":
			batchwiseMatches.y21++
		case "22":
			batchwiseMatches.y22++
		case "23":
			batchwiseMatches.y23++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"totalRegisters":  femaleRegisters + maleRegisters,
		"femaleRegisters": femaleRegisters,
		"maleRegisters":   maleRegisters,
		"batchwiseRegistration": gin.H{
			"y19": registers.y19,
			"y20": registers.y20,
			"y21": registers.y21,
			"y22": registers.y22,
			"y23": registers.y23,
		},
		"totalMatches": numberOfMatches,
		"batchwiseMatches": gin.H{
			"y19": batchwiseMatches.y19,
			"y20": batchwiseMatches.y20,
			"y21": batchwiseMatches.y21,
			"y22": batchwiseMatches.y22,
			"y23": batchwiseMatches.y23,
		},
	})
}
