package controllers

import (
	"errors"
	// "fmt"
	"net/http"
	"time"

	"github.com/Akhilstaar/me-my_encryption/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserFirstLogin(c *gin.Context) {
	// User already authenticated in router.go by gin.HandlerFunc

	// Validate the input format
	info := new(models.TypeUserFirst)
	if err := c.BindJSON(info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}
	
	if info.AuthCode == " " {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already registered."})
		return
	}

	// See U later ;) ...
	user := models.User
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
		Id:    info.Id,
		Pass:  info.PassHash,
		PrivK: info.PrivKey,
		PubK:  info.PubKey,
		AuthC: " ",
		Data:  info.Data,
		Dirty: true,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
		return
	}

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
	if user.Submit == true {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Hearts already sent."})
		return
	}

	if info.ENC1 != "" && info.SHA1 != "" {
		newheart1 := models.SendHeart{
			SHA:            info.SHA1,
			ENC:            info.ENC1,
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

		if err := ReturnClaimedHeart(enc, sha, userID.(string)); err != nil {
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

// need to change the flow a bit.
type HeartClaimError struct {
	Message string
}

func (e HeartClaimError) Error() string {
	return e.Message
}

func ReturnClaimedHeart(enc string, sha string, userId string) error {
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
		SHA: sha,
		ENC: enc,
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

	heartModel := models.SendHeart{}
	verifyheart := Db.Model(&heartModel).Where("sha = ? AND enc = ?", info.SHA, info.Enc).First(&heartModel)
	if verifyheart.Error != nil {
		if errors.Is(verifyheart.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Heart Claim Request."})
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
		Id:   info.Enc,
		SHA:  info.SHA,
		Roll: userID.(string),
	}
	if err := Db.Create(&heartclaim).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// TODO: (RESOLVED) Implement "SendClaimedHeartBack" token logic -- DONE ?
	// generate a token for "SendClaimedHeartBack" which is valid for 10? mins.

	c.JSON(http.StatusAccepted, gin.H{"message": "Heart Claim Success"})
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
		if err := ReturnClaimedHeart(enc, sha, userID.(string)); err != nil {
			c.JSON(http.StatusAccepted, gin.H{"message": "Found invalid Claim Requests. It will be recorded"})
			return
		}
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Congrats !!, we just avoided unexpected event with probability < 1/1000."})
}