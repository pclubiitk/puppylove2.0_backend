package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pclubiitk/puppylove2.0_backend/models"
	"github.com/pclubiitk/puppylove2.0_backend/redisclient"
	"golang.org/x/exp/rand"
)

func UpdateAbout(c *gin.Context) {
	about := new(models.UpdateAbout)
	if err := c.BindJSON(about); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}
	if len(about.About) > 60 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Too long about."})
		return
	}
	userID, _ := c.Get("user_id")
	user := models.User{}

	record := Db.Model((&user)).Where("id = ?", userID).Update("about", about.About)
	if record.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Some Error occured Please try later"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "Update Successful!"})

}

func UpdateInterest(c *gin.Context) {
	interestReq := new(models.UpdateInterest)
	if err := c.BindJSON(interestReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}
	userID, _ := c.Get("user_id")
	user := models.User{}

	// to save our server form very very long tags.
	if len(interestReq.Interests) > 40 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Too long tags."})
		return
	}

	record := Db.Model((&user)).Where("id = ?", userID).Update("interests", interestReq.Interests)
	if record.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Some Error occured Please try later"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "Update Successful!"})
}

func GetAllUsersInfo(c *gin.Context) {
	redisUserAboutMap, err1 := redisclient.RedisClient.HGetAll(redisclient.Ctx, "about_map").Result()
	redisUserinterestMap, err2 := redisclient.RedisClient.HGetAll(redisclient.Ctx, "interests_map").Result()

	// currently we are checking both err1 and err2 to be nil,
	// we can change it to check one at a time.
	// but if we request the data base twice like it will be very heavy task.

	if err1 == nil && err2 == nil && len(redisUserAboutMap) > 0 && len(redisUserinterestMap) > 0 {
		c.JSON(http.StatusOK, gin.H{"about": redisUserAboutMap, "interests": redisUserinterestMap})
		return
	}
	var usersInfo []models.UserInfo
	var userModel models.User

	// later we can, modify it to only return the active users
	fetchUsersInfo := Db.Model(&userModel).Select("id", "about", "interests").Where("dirty = ?", true).Find(&usersInfo)
	if fetchUsersInfo.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Some error occured"})
		return
	}
	aboutMap := make(map[string]string)
	interestsMap := make(map[string]string)
	for _, user := range usersInfo {
		redisclient.RedisClient.HSet(redisclient.Ctx, "about_map", user.Id, string(user.About))
		redisclient.RedisClient.HSet(redisclient.Ctx, "interests_map", user.Id, string(user.Interests))

		// setting expiry time for the keys to 3 hours
		redisclient.RedisClient.Expire(redisclient.Ctx, "about_map", 3*time.Hour)
		redisclient.RedisClient.Expire(redisclient.Ctx, "interests_map", 3*time.Hour)
		aboutMap[user.Id] = string(user.About)
	}
	c.JSON(http.StatusOK, gin.H{"about": aboutMap, "interests": interestsMap})
}

func SuggestRandom(c *gin.Context) {
	var users []models.User
	var userDB models.User

	Db.Model(userDB).Where("dirty = ?", true).Find(&users)

	// Shuffle the users randomly
	rand.Seed(uint64(time.Now().UnixNano()))
	rand.Shuffle(len(users), func(i, j int) {
		users[i], users[j] = users[j], users[i]
	})

	// Select the first 10 users
	var results []string
	for i, user := range users {
		if i >= 10 {
			break
		}
		results = append(results, user.Id)
	}
	c.JSON(http.StatusOK, gin.H{"users": results})
}
