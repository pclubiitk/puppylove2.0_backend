package models

import (
	"gorm.io/gorm"
)

type Heart struct {
	ENC        string `json:"enc" binding:"required"`
	SHA        string `json:"sha" binding:"required"`
	Id_encrypt string `json:"id_encrypt" binding:"required"`
}

type Hearts struct {
	Heart1 Heart `json:"heart1" binding:"required"`
	Heart2 Heart `json:"heart2" binding:"required"`
	Heart3 Heart `json:"heart3" binding:"required"`
	Heart4 Heart `json:"heart4" binding:"required"`
}

type SendHeartVirtual struct {
	Hearts Hearts `json:"hearts" binding:"required"`
}

type SendHeartFirst struct {
	GenderOfSender string             `json:"genderOfSender" binding:"required"`
	ENC1           string             `json:"enc1" binding:"required"`
	SHA1           string             `json:"sha1" binding:"required"`
	ENC2           string             `json:"enc2"`
	SHA2           string             `json:"sha2"`
	ENC3           string             `json:"enc3"`
	SHA3           string             `json:"sha3"`
	ENC4           string             `json:"enc4"`
	SHA4           string             `json:"sha4"`
	ReturnHearts   []VerifyHeartClaim `json:"returnhearts"`
}

type VerifyHeartClaim struct {
	Enc            string `json:"enc" binding:"required"`
	SHA            string `json:"sha" binding:"required"`
	GenderOfSender string `json:"genderOfSender" binding:"required"`
}

type FetchHeartsFirst struct {
	Enc            string `json:"enc"`
	GenderOfSender string `json:"genderOfSender"`
}
type FetchReturnedHearts struct {
	SHA string `json:"sha"`
	Enc string `json:"enc"`
}
type SentHeartsDecoded struct {
	DecodedHearts []FetchHeartsFirst `json:"decodedHearts" binding:"required"`
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
		Id   string `json:"enc" bson:"enc" gorm:"unique"`
		SHA  string `json:"sha" bson:"sha" gorm:"unique"`
		Roll string `json:"roll"`
	}
)

// --------- Returning Heart Below ---------

type UserReturnHearts struct {
	ReturnHearts []UserReturnHeart `json:"returnhearts" binding:"required"`
}

type UserReturnHeart struct {
	ENC string `json:"enc" binding:"required" gorm:"unique"`
	SHA string `json:"sha" binding:"required" gorm:"unique"`
}

type (
	ReturnHearts struct {
		gorm.Model
		SHA string `json:"sha" bson:"sha"`
		ENC string `json:"enc" bson:"enc" gorm:"unique"`
	}
)

type FetchReturnHeart struct {
	ENC string `json:"enc" binding:"required" gorm:"unique"`
}

type (
	MatchTable struct {
		gorm.Model
		Roll1 string `json:"roll1" bson:"roll1"`
		Roll2 string `json:"roll2" bson:"roll2"`
	}
)
