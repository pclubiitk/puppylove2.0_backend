package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func CheckGoogleCaptcha(response string) bool {
	var googleCaptcha string = os.Getenv("CAPTCHA_SECRET")
	req, _ := http.NewRequest("POST", "https://www.google.com/recaptcha/api/siteverify", nil)
	q := req.URL.Query()
	q.Add("secret", googleCaptcha)
	q.Add("response", response)
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	var googleResponse map[string]interface{}
	resp, err := client.Do(req)
	// fmt.Println(resp)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &googleResponse)
	fmt.Println(googleResponse["success"])
	return googleResponse["success"].(bool)
}

func Captchacheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		captcha := c.GetHeader("g-recaptcha-response")
		// function working properly, no need for this now.
		// fmt.Println(captcha)
		human := CheckGoogleCaptcha(captcha)
		if !human {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Captcha."})
			c.Abort()
			return
		}
		c.Next()
	}
}

func AuthenticateAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var jwtSigningKey = os.Getenv("USER_JWT_SIGNING_KEY")
		authCookie, err := c.Cookie("Authorization")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization cookie missing"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(authCookie, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
				c.Abort()
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSigningKey), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := claims["user_id"].(string)
			if userID != os.Getenv("ADMIN_ID") {
				unauthorizedMessage := fmt.Sprintf("Unauthorized Login attempt by %s, it will be recorded.", userID)
				log.Println(unauthorizedMessage)
				c.JSON(http.StatusUnauthorized, gin.H{"error": unauthorizedMessage})
				c.Abort()
			}
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
		}
	}
}

func AuthenticateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var jwtSigningKey = os.Getenv("USER_JWT_SIGNING_KEY")
		authCookie, err := c.Cookie("Authorization")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization cookie missing"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(authCookie, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
				c.Abort()
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSigningKey), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := claims["user_id"].(string)
			if userID == os.Getenv("ADMIN_ID") {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "You Serious ?? , login using Postman or CLI, the frontend is only for normal users."})
				c.Abort()
			}
			c.Set("user_id", userID)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
		}
	}
}

func AuthenticateUserHeartclaim() gin.HandlerFunc {
	return func(c *gin.Context) {
		var heartjwtSigningKey = os.Getenv("HEART_JWT_SIGNING_KEY")
		useridfromJWT, _ := c.Get("user_id")
		heartCookie, err := c.Cookie("HeartBack")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Late HeartClaim cookie missing"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(heartCookie, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token Signing Algo"})
				c.Abort()
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(heartjwtSigningKey), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Heart Token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			HeartBackUserID := claims["user_id"]
			verified := claims["verified"]
			// Added this for safety, if someone runs the server with same keys for both.
			// HeartjwtSigningKey, jwtSigningKey should be different.
			if verified != "Absolutely" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Forged Claim Token, same as Authorization."})
				c.Abort()
				return
			}
			fmt.Println(HeartBackUserID)
			if useridfromJWT != HeartBackUserID {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Forged Claim Token, invalid Used_Id"})
				c.Abort()
				return
			}
			c.Set("HeartBackUserID", HeartBackUserID)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
		}
	}
}

func AdminPermit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if permit {
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not Permitted by Admin"})
			c.Abort()
		}
	}
}
