package controllers

import (
	"errors"
	"net/http"
	"strings"
	"time"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pclubiitk/puppylove2.0_backend/mail"
	"github.com/pclubiitk/puppylove2.0_backend/models"
	"github.com/pclubiitk/puppylove2.0_backend/utils"
	"gorm.io/gorm"
	"github.com/pclubiitk/puppylove2.0_backend/redisclient"
)

func FetchPublicKeys(c *gin.Context) {
	publicKeysMap, err := redisclient.RedisClient.HGetAll(redisclient.Ctx, "public_keys").Result()
	if err == nil && len(publicKeysMap) > 0 {
		c.JSON(http.StatusOK, publicKeysMap)
		return
	}
	var publicKeys []models.UserPublicKey
	var userModel models.User
	fetchPublicKey := Db.Model(&userModel).Select("id, pub_k").Find(&publicKeys)
	if fetchPublicKey.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Some Error Occurred"})
		return
	}
	responseMap := make(map[string]string)
	for _, key := range publicKeys {
		redisclient.RedisClient.HSet(redisclient.Ctx, "public_keys", key.Id, key.PubK)
		responseMap[key.Id] = key.PubK
	}
	// redisclient.ViewRedis()
	c.JSON(http.StatusOK, responseMap)
}

func FetchReturnHearts(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in context"})
		return
	}

	type Time struct {
		ReturnHeartsTimestamp string `json:"return_hearts_timestamp"`
	}
	var userTimestamp Time
	if err := Db.Model(&models.User{}).
		Select("return_hearts_timestamp").
		Where("id = ?", userID).
		Scan(&userTimestamp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user timestamp"})
		}
		return
	}
	layout := "2006-01-02T15:04:05.999999999Z07:00" // RFC 3339 format
	timestamp, err := time.Parse(layout, userTimestamp.ReturnHeartsTimestamp)
	if err != nil {
		fmt.Println("Failed to parse timestamp:", userTimestamp.ReturnHeartsTimestamp)
		fmt.Println("Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp format"})
		return
	}
	var returnedHeart models.ReturnHearts
	var returnedHearts []models.FetchReturnedHearts
	if err := Db.Model(&returnedHeart).
		Select("enc").
		Where("created_at > ?", timestamp).
		Find(&returnedHearts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch returned hearts"})
		return
	}
	newTimestamp := time.Now().UTC().Add(-1 * time.Minute).Format(time.RFC3339)
	if err := Db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("return_hearts_timestamp", newTimestamp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user timestamp"})
		return
	}

	c.JSON(http.StatusOK, returnedHearts)
}

func FetchHearts(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in context"})
		return
	}
	type Time struct {
		SendHeartsTimestamp string `json:"send_hearts_timestamp"`
	}
	var userTimestamp Time

	if err := Db.Model(&models.User{}).
		Select("send_hearts_timestamp").
		Where("id = ?", userID).
		Scan(&userTimestamp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user timestamp"})
		}
		return
	}
	layout := "2006-01-02T15:04:05.999999999Z07:00" // RFC 3339 format
	timestamp, err := time.Parse(layout, userTimestamp.SendHeartsTimestamp)
	if err != nil {
		fmt.Println("Failed to parse timestamp:", userTimestamp.SendHeartsTimestamp)
		fmt.Println("Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp format"})
		return
	}
	var heart models.SendHeart
	var hearts []models.FetchHeartsFirst

	if err := Db.Model(&heart).
		Select("enc, gender_of_sender").
		Where("created_at > ?", timestamp).
		Find(&hearts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hearts"})
		return
	}
	newTimestamp := time.Now().UTC().Add(-1*time.Minute).Format(time.RFC3339)
	if err := Db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("send_hearts_timestamp", newTimestamp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user timestamp"})
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
