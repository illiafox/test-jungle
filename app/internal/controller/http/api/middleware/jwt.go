package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	jwt2 "jungle-test/internal/controller/http/api/jwt"
)

type JwtMiddleware struct {
	jc *jwt2.JwtConfigurator
}

func NewAuthMiddleware(jc *jwt2.JwtConfigurator) *JwtMiddleware {
	return &JwtMiddleware{jc: jc}
}

type (
	userClaimsKey struct{}
)

var (
	ErrNoAuthorizationHeader    = fiber.NewError(fiber.StatusUnauthorized, "'Authorization' header wasn't provided")
	ErrWrongAuthorizationFormat = fiber.NewError(fiber.StatusUnauthorized, "wrong authorization format")
)

func (m JwtMiddleware) AccessTokenMiddleware(ctx *fiber.Ctx) error {
	authorization := ctx.Get(fiber.HeaderAuthorization)
	if authorization == "" {
		return ErrNoAuthorizationHeader
	}

	parts := strings.Split(authorization, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ErrWrongAuthorizationFormat
	}

	userClaims, err := m.jc.VerifyAccessToken(parts[1])
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	ctx.Locals(userClaimsKey{}, userClaims)
	return ctx.Next()
}

func GetUserClaims(ctx *fiber.Ctx) jwt2.UserClaims {
	return ctx.Locals(userClaimsKey{}).(jwt2.UserClaims)
}
