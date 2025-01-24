package models

import (
	"gorm.io/gorm"
)

var PublishMatches = false

type (
	// User represents the structure of our resource
	User struct {
		gorm.Model
		Id       string `json:"_id" bson:"_id" gorm:"unique"`
		Name     string `json:"name" bson:"name"`
		Email    string `json:"email" bson:"email" gorm:"unique"`
		Gender   string `json:"gender" bson:"gender"`
		Pass     string `json:"passHash" bson:"passHash"`
		PubK     string `json:"pubKey" bson:"pubKey"`
		PrivK    string `json:"privKey" bson:"privKey"`
		AuthC    string `json:"authCode" bson:"authCode"`
		Data     string `json:"data" bson:"data"`
		Claims   string `json:"claims" bson:"claims"`
		Submit   bool   `json:"submitted" bson:"submitted"`
		Matches  string `json:"matches" bson:"matches"`
		Dirty    bool   `json:"dirty" bson:"dirty"`
		Publish  bool   `json:"publish" bson:"publish"`
		Code     string `json:"code" bson:"code"`
		About    string `json:"about" bson:"about"`
		Intrests string `json:"intrests" bson:"intrests"`
	}
)
type UserPublicKey struct {
	gorm.Model
	Id   string `json:"_id" bson:"_id" gorm:"unique"`
	PubK string `json:"pubKey" bson:"pubKey"`
}

type AddNewUser struct {
	TypeUserNew []TypeUserNew `json:"newuser" binding:"required"`
}

type TypeUserNew struct {
	Id     string `json:"roll" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Email  string `json:"email" binding:"required"`
	Gender string `json:"gender" binding:"required"`
}

type TypeUserFirst struct {
	Id       string `json:"roll" binding:"required"`
	AuthCode string `json:"authCode" binding:"required"`
	PassHash string `json:"passHash" binding:"required"`
	PubKey   string `json:"pubKey" binding:"required"`
	PrivKey  string `json:"privKey" binding:"required"`
	Data     string `json:"data" binding:"required"`
}

type UserLogin struct {
	Id   string `json:"_id" binding:"required"`
	Pass string `json:"passHash" binding:"required"`
}

type RecoveryCodeReq struct {
	Pass string `json:"passHash" binding:"required"`
	Code string `json:"code" binding:"required"`
}

type RetrivePassReq struct {
	Id string `json:"_id" binding:"required"`
}

type UpdateAbout struct {
	About string `json:"about" binding:"required"`
}

type UpdateIntrest struct {
	Intrests string `json:"intrests"`
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
	Dirty bool   `json:"dirty" binding:"required"`
}

var StatsFlag = true
var FemaleRegisters = 0
var MaleRegisters = 0
var NumberOfMatches = 0
var RegisterMap = make(map[string]int)
var MatchMap = make(map[string]int)
