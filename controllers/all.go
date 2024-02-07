package controllers

import (
	"errors"
	"net/http"
	"strings"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/pclubiitk/puppylove2.0_backend/mail"
	"github.com/pclubiitk/puppylove2.0_backend/models"
	"github.com/pclubiitk/puppylove2.0_backend/utils"
	"github.com/dpapathanasiou/go-recaptcha"
	"fmt"
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

func VerifyCaptcha(clientIP, token string) error {
	csec := os.Getenv("CAPTCHA_SECRET")
	recaptcha.Init(csec)

	res,err := recaptcha.Confirm(clientIP, token)
	// fmt.Println(res)
	// fmt.Println(token)
	if res == false {
		return err
	}

	return nil
}

func UserMail(c *gin.Context) {
	id := c.Param("id")

	recap := new(models.Captcha)
	if err := c.BindJSON(recap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Captcha data format."})
		return
	}
	clientIP := c.ClientIP()
	// fmt.Println(clientIP)
	// fmt.Println(recap.recaptchatoken)

	err := VerifyCaptcha(clientIP, recap.recaptchatoken)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Captcha token."})
		return
	}

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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Something went wrong, Please try again."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Auth. code sent successfully !!"})
}

func GetStats(c *gin.Context) {
	if !models.PublishMatches {
		c.JSON(http.StatusOK, gin.H{"msg": "Stats Not yet published"})
		return
	}
	if models.StatsFlag {
		models.StatsFlag = false
		var userdb models.User
		var users []models.User

		records := Db.Model(&userdb).Where("dirty = ?", true).Find(&users)
		if records.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Fetching Stats"})
			return
		}

		for _, user := range users {
			// log.Print(user.Id)
			if len(user.Id) >= 2 && user.Dirty {
				if user.Gender == "M" {
					models.MaleRegisters++
				} else {
					models.FemaleRegisters++
				}
				models.RegisterMap["y"+user.Id[0:2]]++
				if user.Matches == "" {
					continue
				}
				var matchCount = len(strings.Split(user.Matches, ","))
				if matchCount != 0 {
					models.NumberOfMatches += matchCount
					var myMatches = strings.Split(user.Matches, ",")
					for _, t := range myMatches {
						if len(t) >= 2 {
							models.MatchMap["y"+t[0:2]]++
						}
					}
				}
			}
			// log.Print("Done")
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"totalRegisters":        models.MaleRegisters + models.FemaleRegisters,
		"femaleRegisters":       models.FemaleRegisters,
		"maleRegisters":         models.MaleRegisters,
		"batchwiseRegistration": models.RegisterMap,
		"totalMatches":          models.NumberOfMatches,
		"batchwiseMatches":      models.MatchMap})
}
