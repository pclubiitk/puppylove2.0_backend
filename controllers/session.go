package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pclubiitk/puppylove2.0_backend/models"
	"gorm.io/gorm"
)

func UserLogin(c *gin.Context) {
	info := new(models.UserLogin)
	if err := c.BindJSON(info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}

	loginmodel := models.User{}
	verifyuser := Db.Model(&loginmodel).Where("id = ? AND pass = ?", info.Id, info.Pass).First(&loginmodel)
	if verifyuser.Error != nil {
		fmt.Println(verifyuser.Error)
		if errors.Is(verifyuser.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Login Request."})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": verifyuser.Error.Error()})
			return
		}
	}

	token, err := generateJWTToken(info.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
		return
	}
	expirationTime := time.Now().Add(time.Hour * 24)
	cookie := &http.Cookie{
		Name:    "Authorization",
		Value:   token,
		Expires: expirationTime,
		Path:    "/",
		Domain:  os.Getenv("DOMAIN"),
		// For Http
		// HttpOnly: true,
		// Secure:   false, // Set this to true if you're using HTTPS, false for HTTP
		// SameSite: http.SameSiteStrictMode,
		// For Https
		HttpOnly: false,
		Secure:   true, // Set this to true if you're using HTTPS, false for HTTP
		SameSite: http.SameSiteNoneMode,
	}
	// at the time of login, we set the auth cookie and send back the pub and priv keys
	http.SetCookie(c.Writer, cookie)
	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully !!", "pvtKey_Enc": loginmodel.PrivK, "pubKey": loginmodel.PubK})
}

func AddRecovery(c *gin.Context) {
	// Validate the input format
	data := new(models.RecoveryCodeReq)
	if err := c.BindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}
	user := models.User{}

	userId, _ := c.Get("user_id")
	// find the record, auth insures the user is valid
	result := Db.Model(&user).Where("id = ? AND pass = ? ", userId, data.Pass).First(&user).Update("Code", data.Code)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Some Error occured Please try later"})
	} else {
		c.JSON(http.StatusAccepted, gin.H{"message": "Recovery Code Added Sucessfully!"})
	}
}

func RetrivePass(c *gin.Context) {
	req := new(models.RetrivePassReq)
	if err := c.BindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}
	user := models.User{}
	record := Db.Model(&user).Where("id = ?", req.Id).First(&user)
	if record.Error != nil {
		if errors.Is(record.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{"error": "User Does Not Exist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Some Error occured Please try later"})
		return
	} else {
		passCode := user.Code
		if passCode == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Didn't Registered With Recovery Codes"})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"message": "Successfully retrived code", "code": passCode})
		return
	}
}

func GetUserData(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user := models.User{}
	result := Db.Model(&user).Where("id = ?", userID).First(&user)
	if result.Error != nil {
		// here we assume after authentication user id will always be present, if not it is internal server error
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data retrieved successfully !!", "id": userID, "data": user.Data, "gender": user.Gender, "submit": user.Submit, "claims": user.Claims, "permit": permit, "publish": user.Publish, "about": user.About, "intrest": user.Intrests})
}

func UserLogout(c *gin.Context) {
	cookie := &http.Cookie{
		Name:    "Authorization",
		Value:   "",
		Expires: time.Unix(0, 0), // Expire the cookie immediately
		// MaxAge:  -1,
		Path:     "/",
		Domain:   os.Getenv("DOMAIN"),
		// HttpOnly: true,
		// Secure:   false,
		// SameSite: http.SameSiteStrictMode,
		// For Htpps
		HttpOnly: false,
		Secure:   true, // Set this to true if you're using HTTPS, false for HTTP
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(c.Writer, cookie)

	c.JSON(http.StatusOK, gin.H{
		"message": "User logged out successfully.",
	})
}

type AuthClaims struct {
	User_id string `json:"user_id"`
	jwt.StandardClaims
}
type HeartClaims struct {
	User_id  string `json:"user_id"`
	Verified string `json:"verified"`
	jwt.StandardClaims
}

func generateJWTToken(userID string) (string, error) {
	var jwtSigningKey = os.Getenv("USER_JWT_SIGNING_KEY")
	claims := AuthClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
			IssuedAt:  jwt.TimeFunc().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSigningKey))
	return tokenString, err
}

func generateJWTTokenForHeartBack(userID string) (string, error) {
	var heartjwtSigningKey = os.Getenv("HEART_JWT_SIGNING_KEY")
	verified := "Absolutely"
	claims := HeartClaims{
		userID,
		verified,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour / 3).Unix(),
			IssuedAt:  jwt.TimeFunc().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(heartjwtSigningKey))
	return tokenString, err
}
