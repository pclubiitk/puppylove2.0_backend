package controllers

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pclubiitk/puppylove2.0_backend/models"
	"github.com/pclubiitk/puppylove2.0_backend/redisclient"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	// "strings"
	"time"
)

func UserFirstLogin(c *gin.Context) {
	// User already authenticated in router.go by gin.HandlerFunc

	// Validate the input format
	info := new(models.TypeUserFirst)
	if err := c.BindJSON(info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}

	tempU := models.MailData{}
	tempUser := models.User{}
	tempRecord := Db.Model(&tempUser).Where("id = ?", info.Id).First(&tempU)
	if tempRecord.Error != nil {
		if errors.Is(tempRecord.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{"error": "User not found !!"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
			return
		}
	}
	if tempU.Dirty {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "User already registered"})
		return
	}

	// if info.AuthCode == " " {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "User already registered."})
	// 	return
	// }

	// See U later ;) ...
	user := models.User{}
	publicK := Db.Model(&user).Where("pub_k = ?", info.PubKey).First(&user)
	if publicK.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter another public key !!"})
		return
	}

	record := Db.Model(&user).Where("id = ? AND auth_c = ?", info.Id, info.AuthCode).First(&user)
	if record.Error != nil {
		if errors.Is(record.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Incorrect AuthCode entered !!"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
			return
		}
	}

	// var newuser models.User
	if err := record.Updates(models.User{
		Id:     info.Id,
		Pass:   info.PassHash,
		PubK:   info.PubKey,
		PrivK:  info.PrivKey,
		AuthC:  " ",
		Data:   info.Data,
		Claims: "",
		Dirty:  true,
		Code:   "",
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
		return
	}
	redisclient.RedisClient.HSet(redisclient.Ctx, "public_keys", info.Id, info.PubKey)
	c.JSON(http.StatusCreated, gin.H{"message": "User Created Successfully."})
}

func SendHeart(c *gin.Context) {
	// User already authenticated in router.go by gin.HandlerFunc

	info := new(models.SendHeartFirst)
	if err := c.BindJSON(info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}

	userID, _ := c.Get("user_id")
	var user models.User
	record := Db.Model(&user).Where("id = ?", userID).First(&user)
	if record.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
		return
	}
	if user.Submit {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Hearts already sent."})
		return
	}

	if info.ENC1 != "" && info.SHA1 != "" {
		newheart1 := models.SendHeart{
			SHA:            info.SHA1,
			ENC:            info.ENC1,
			SONG_ENC:       info.SONG1,
			GenderOfSender: info.GenderOfSender,
		}
		if err := Db.Create(&newheart1).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
	}

	if err := record.Updates(models.User{
		Submit: true,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
		return
	}

	if info.ENC2 != "" && info.SHA2 != "" {
		newheart2 := models.SendHeart{
			SHA:            info.SHA2,
			ENC:            info.ENC2,
			SONG_ENC:       info.SONG2,
			GenderOfSender: info.GenderOfSender,
		}
		if err := Db.Create(&newheart2).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error in submitting heart 1": err})
			return
		}
	}

	if info.ENC3 != "" && info.SHA3 != "" {
		newheart3 := models.SendHeart{
			SHA:            info.SHA3,
			ENC:            info.ENC3,
			SONG_ENC:       info.SONG3,
			GenderOfSender: info.GenderOfSender,
		}
		if err := Db.Create(&newheart3).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error in submitting heart 2": err})
			return
		}
	}

	if info.ENC4 != "" && info.SHA4 != "" {
		newheart4 := models.SendHeart{
			SHA:            info.SHA4,
			ENC:            info.ENC4,
			SONG_ENC:       info.SONG4,
			GenderOfSender: info.GenderOfSender,
		}
		if err := Db.Create(&newheart4).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error in submitting heart 3": err})
			return
		}
	}

	for _, heart := range info.ReturnHearts {
		enc := heart.Enc
		sha := heart.SHA
		song := heart.SONG_ENC

		if err := ReturnClaimedHeart(enc, sha, song, userID.(string)); err != nil {
			c.JSON(http.StatusAccepted, gin.H{"message": "Hearts Sent Successfully !!, but found invalid Claim Requests. It will be recorded"})
			return
		}
	}

	token, err := generateJWTTokenForHeartBack(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
		return
	}
	expirationTime := time.Now().Add(time.Hour * 24)
	cookie := &http.Cookie{
		Name:     "HeartBack",
		Value:    token,
		Expires:  expirationTime,
		Path:     "/",
		Domain:   "localhost",
		HttpOnly: true,
		Secure:   false, // Set this to true if you're using HTTPS, false for HTTP
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(c.Writer, cookie)

	c.JSON(http.StatusAccepted, gin.H{"message": "Hearts Sent Successfully !!"})
}

func SendHeartVirtual(c *gin.Context) {
	info := new(models.SendHeartVirtual)
	if err := c.BindJSON(info); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong Format"})
		return
	}

	userID, _ := c.Get("user_id")
	var user models.User
	record := Db.Model(&user).Where("id = ?", userID).First(&user)
	if record.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User does not exist."})
		return
	}

	if user.Submit {
		c.JSON(http.StatusOK, gin.H{"error": "Hearts already sent."})
		return
	}

	jsonData, err := json.Marshal(info.Hearts)
	if err != nil {
		return
	}

	if err := record.Updates(models.User{
		Data: string(jsonData),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update Data field of User."})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Virtual Hearts Sent Successfully !!"})
}

// need to change the flow a bit.
type HeartClaimError struct {
	Message string
}

func (e HeartClaimError) Error() string {
	return e.Message
}

func ReturnClaimedHeart(enc string, sha string, song string, userId string) error {
	heartModel := models.HeartClaims{}
	if enc == "" || sha == "" {
		return nil
	}
	verifyheart := Db.Model(&heartModel).Where("sha = ? AND roll = ?", sha, userId).First(&heartModel)
	if verifyheart.Error != nil {
		if errors.Is(verifyheart.Error, gorm.ErrRecordNotFound) {
			return HeartClaimError{Message: "Unauthorized Heart Claim attempt, it will be recorded."}
		} else {
			return HeartClaimError{Message: "verifyheart.Error.Error()"}
		}
	}

	heartclaim := models.ReturnHearts{
		SHA:      sha,
		ENC:      enc,
		SONG_ENC: song,
	}
	if err := Db.Create(&heartclaim).Error; err != nil {
		return HeartClaimError{Message: "Something Unexpected Occurred while adding the heart claim."}
	}

	return nil
}

func HeartClaim(c *gin.Context) {

	info := new(models.VerifyHeartClaim)
	if err := c.BindJSON(info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}

	claim_status := "true"
	heartModel := models.SendHeart{}
	verifyheart := Db.Model(&heartModel).Where("sha = ? AND enc = ?", info.SHA, info.Enc).First(&heartModel)
	if verifyheart.Error != nil {
		if errors.Is(verifyheart.Error, gorm.ErrRecordNotFound) {
			// c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Heart Claim Request."})
			claim_status = "false"
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": verifyheart.Error.Error()})
		}
		return
	}
	// If the db has record of sha and enc then remove it from the record and add the sha, enc to userId
	if err := Db.Model(&heartModel).Where("sha = ? AND enc = ?", info.SHA, info.Enc).Unscoped().Delete(&heartModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	heartclaim := models.HeartClaims{
		Id:       info.Enc,
		SHA:      info.SHA,
		Roll:     userID.(string),
		SONG_ENC: heartModel.SONG_ENC,
	}
	if err := Db.Create(&heartclaim).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var user models.User
	record := Db.Model(&user).Where("id = ?", userID).First(&user)
	if record.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User does not exist."})
		return
	}

	jsonClaim, err := json.Marshal(info)
	if err != nil {
		return
	}

	newClaim := string(jsonClaim)
	// Encoding '+' present claim so that the claim can be concatenated with '+' and later retrieved
	newClaim_enc := url.QueryEscape(newClaim)

	if user.Claims == "" {
		user.Claims = newClaim_enc
	} else {
		user.Claims = user.Claims + "+" + newClaim_enc
	}

	if err := Db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update Claims field of User."})
		return
	}

	// TODO: (RESOLVED) Implement "SendClaimedHeartBack" token logic -- DONE ?
	// generate a token for "SendClaimedHeartBack" which is valid for 10? mins.

	c.JSON(http.StatusAccepted, gin.H{"message": "Heart Claim Success", "claim_status": claim_status})
}

// TODO: (RESOLVED) Current issue is that if the user changes the enc of the claimed hash(which is very timeconsuming btw ;), there is no way to verify here. -- DONE
// Why not just add a time window of 10? mins in which the heartback can be accessed.
// So, what are the odds that user gets a heart within 10 mins of submitting its hearts ?.
// Even if the user gets it, what are the odds that user will be able to Intercept the request and make a claim with "enc" which is encoded with pub key of user's 5th choice ?
func ReturnClaimedHeartLate(c *gin.Context) {
	// TODO: Modify this function to handle multiple concatenated json inputs

	info := new(models.UserReturnHearts)
	if err := c.BindJSON(info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}

	userID, _ := c.Get("user_id")
	for _, heart := range info.ReturnHearts {
		enc := heart.ENC
		sha := heart.SHA
		song := heart.SONG_ENC
		if err := ReturnClaimedHeart(enc, sha, song, userID.(string)); err != nil {
			c.JSON(http.StatusAccepted, gin.H{"message": "Found invalid Claim Requests. It will be recorded"})
			return
		}
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Congrats !!, we just avoided unexpected event with probability < 1/1000."})
}

func Publish(c *gin.Context) {
	if models.PublishMatches {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Results Published"})
		return
	}
	userID, _ := c.Get("user_id")
	var user models.User
	record := Db.Model(&user).Where("id = ?", userID).Update("Publish", true)
	if record.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
		return
	}
}

func GetActiveUsers(c *gin.Context) {
	var users []models.User
	var userDB models.User
	Db.Model(userDB).Find(&users)
	var results []string
	for _, user := range users {
		if user.Dirty {
			results = append(results, user.Id)
		}
	}
	c.JSON(http.StatusOK, gin.H{"users": results})
}

/*
This function would verify heart claims from returned table and would take care match logic adding matched rollno to matching table
*/
func VerifyReturnHeart(c *gin.Context) {
	info := new(models.VerifyReturnHeartClaim)
	if err := c.BindJSON(info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}
	h := sha256.New()
	h.Write([]byte(info.Secret))
	bs := h.Sum(nil)
	hash := fmt.Sprintf("%x", bs)
	heartModel := models.ReturnHearts{}
	verifyheart := Db.Model(&heartModel).Where("sha = ? AND enc = ?", hash, info.Enc).First(&heartModel)
	if verifyheart.Error != nil {
		if errors.Is(verifyheart.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Heart Claim Request."})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": verifyheart.Error.Error()})
		}
		return
	}
	if err := Db.Model(&heartModel).Where("sha = ? AND enc = ?", hash, info.Enc).Unscoped().Delete(&heartModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var heartClaim models.HeartClaims
	Db.Model(heartClaim).Where("sha = ?", hash).First(&heartClaim)
	userID, _ := c.Get("user_id")
	// roll1 := userID.(string)
	// roll2 := heartClaim.Roll

	// userdb := models.User{}
	// Db.Model(&userdb).Where("id = ?", roll1).First(&userdb)
	// userdb.Matches = userdb.Matches + "," + roll2
	// Db.Save(&userdb)

	// userdb2 := models.User{}
	// Db.Model(&userdb2).Where("id = ?", roll2).First(&userdb2)
	// userdb2.Matches = userdb2.Matches + "," + roll1
	// Db.Save(&userdb2)

	// temp1, _ := strconv.Atoi(userID.(string))
	// temp2, _ := strconv.Atoi(heartClaim.Roll)

	// if temp1 < temp2 {
	// 	returnHeartClaim := models.MatchTable{
	// 		Roll1: userID.(string),
	// 		Roll2: heartClaim.Roll,
	// 	}
	// 	if err := Db.Create(&returnHeartClaim).Error; err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	// 		return
	// 	}
	// } else if temp2 < temp1 {
	// 	returnHeartClaim := models.MatchTable{
	// 		Roll2: userID.(string),
	// 		Roll1: heartClaim.Roll,
	// 	}
	// 	if err := Db.Create(&returnHeartClaim).Error; err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	// 		return
	// 	}
	// }

	var existingMatch models.MatchTable
	if err := Db.Where("roll1 = ? AND roll2 = ?", heartClaim.Roll, userID.(string)).First(&existingMatch).Error; err == nil {
		// Entry already exists, so return without inserting
		c.JSON(http.StatusAccepted, gin.H{"message": "Match Already Done from other side"})
		return
	}
	returnHeartClaim := models.MatchTable{
		Roll1:  userID.(string),
		Roll2:  heartClaim.Roll,
		SONG12: heartClaim.SONG_ENC,
		SONG21: heartModel.SONG_ENC,
	}
	if err := Db.Create(&returnHeartClaim).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Heart Claim Success"})
}

func MatchesHandler(c *gin.Context) {
	if models.PublishMatches {
		userID, _ := c.Get("user_id")
		var user models.User

		if err := Db.Model(&user).Where("id = ?", userID).First(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		if !user.Publish {
			c.JSON(http.StatusOK, gin.H{"msg": "You chose not to publish results"})
			return
		}

		// Unmarshal Matches field for response
		var matches map[string]string
		if len(user.Matches) > 0 {
			_ = json.Unmarshal(user.Matches, &matches)
		} else {
			matches = make(map[string]string)
		}

		c.JSON(http.StatusOK, gin.H{"matches": matches})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Matches not yet published"})
}
