package models

import (
	"gorm.io/gorm"
)

type (
	// User represents the structure of our resource
	User struct {
		gorm.Model
		Id      string `json:"_id" bson:"_id" gorm:"unique"`
		Name    string `json:"name" bson:"name"`
		Email   string `json:"email" bson:"email" gorm:"unique"`
		Gender  string `json:"gender" bson:"gender"`
		Pass    string `json:"passHash" bson:"passHash"`
		PubK    string `json:"pubKey" bson:"pubKey"`
		AuthC   string `json:"authCode" bson:"authCode"`
		Data    string `json:"data" bson:"data"`
		Submit  bool   `json:"submitted" bson:"submitted"`
		Matches string `json:"matches" bson:"matches"`
		Dirty   bool   `json:"dirty" bson:"dirty"`
	}
)

type AddNewUser struct {
	TypeUserNew []TypeUserNew `json:"newuser" binding:"required"`
}

type TypeUserNew struct {
	Id       string `json:"roll" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Gender   string `json:"gender" binding:"required"`
	PassHash string `json:"passHash" binding:"required"`
}

type TypeUserFirst struct {
	Id       string `json:"roll" binding:"required"`
	AuthCode string `json:"authCode" binding:"required"`
	PassHash string `json:"passHash" binding:"required"`
	PubKey   string `json:"pubKey" binding:"required"`
	Data     string `json:"data" binding:"required"`
}

type UserLogin struct {
	Id   string `json:"_id" binding:"required"`
	Pass string `json:"passHash" binding:"required"`
}

// w'll change it later (maybee..)
type AdminLogin struct {
	Id   string `json:"id" binding:"required"`
	Pass string `json:"pass" binding:"required"`
}

type MailData struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	AuthC string `json:"authCode" binding:"required"`
}
