package models

import (
	"gorm.io/gorm"
)

type SendHeartFirst struct {
	GenderOfSender string `json:"genderOfSender" binding:"required"`
	ENC1           string `json:"enc1" binding:"required"`
	SHA1           string `json:"sha1" binding:"required"`
	ENC2           string `json:"enc2"`
	SHA2           string `json:"sha2"`
	ENC3           string `json:"enc3"`
	SHA3           string `json:"sha3"`
	ENC4           string `json:"enc4"`
	SHA4           string `json:"sha4"`
	ReturnHearts   []VerifyHeartClaim `json:"returnhearts"`
}

type VerifyHeartClaim struct {
	Enc string `json:"enc" binding:"required"`
	SHA string `json:"sha" binding:"required"`
}

type FetchHeartsFirst struct {
	Enc            string `json:"enc"`
	GenderOfSender string `json:"genderOfSender"`
}

// gorm.Model represents the structure of our resource in db
type (
	SendHeart struct {
		gorm.Model
		SHA            string `json:"sha" bson:"sha" gorm:"unique"`
		ENC            string `json:"enc" bson:"enc" gorm:"unique"`
		GenderOfSender string `json:"genderOfSender" bson:"gender"`
	}
)

type (
	HeartClaims struct {
		gorm.Model
		Id string `json:"enc" bson:"enc" gorm:"unique"`
		SHA string `json:"sha" bson:"sha" gorm:"unique"`
		Roll  string `json:"roll"`
	}
)

// --------- Returning Heart Below ---------

type UserReturnHearts struct {
	ReturnHearts   []UserReturnHeart `json:"returnhearts" binding:"required"`
}

type UserReturnHeart struct {
	ENC string `json:"enc" binding:"required" gorm:"unique"`
	SHA string `json:"sha" binding:"required" gorm:"unique"`
}

type (
	ReturnHearts struct {
		gorm.Model
		SHA string `json:"sha" bson:"sha" gorm:"unique"`
		ENC string `json:"enc" bson:"enc" gorm:"unique"`
	}
)

type FetchReturnHeart struct {
	ENC string `json:"enc" binding:"required" gorm:"unique"`
}