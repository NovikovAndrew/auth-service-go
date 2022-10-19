package data

import "time"

type User struct {
	ID         string    `json:"id" sql:"id"`
	Email      string    `json:"email" sql:"email"`
	Password   string    `json:"password" sql:"password"`
	Username   string    `json:"username" sql:"username"`
	TokenHash  string    `json:"tokenhash" sql:"tokenhash"`
	IsVerified bool      `json:"isverified" sql:"isverified"`
	CreatedAt  time.Time `json:"createdat" sql:"createdat"`
	UpdatedAt  time.Time `json:"updatedAt" sql:"updatedat"`
}

type VerificationDatatype int

type VerificationData struct {
	Email     string               `json:"email" validate:"required" sql:"email"`
	Code      string               `json:"code" sql:"code"`
	ExpiresAt time.Time            `json:"expiresat" sql:"expiresat"`
	Type      VerificationDatatype `json:"type" sql:"type"`
}
