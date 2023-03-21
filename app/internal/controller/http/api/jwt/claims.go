package jwt

import (
	"github.com/google/uuid"
	"github.com/kataras/jwt"
)

type UserClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string
}

type AccessToken struct {
	jwt.Claims
	UserClaims
}
