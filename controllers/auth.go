package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	// "os"
	// "time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthenticateAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var jwtSigningKey = os.Getenv("UserjwtSigningKey")
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
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
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
			if userID != os.Getenv("AdminId") {
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
		var jwtSigningKey = os.Getenv("UserjwtSigningKey")
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
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
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
			if userID == os.Getenv("AdminId") {
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
		var heartjwtSigningKey = os.Getenv("HeartjwtSigningKey")
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
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
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
