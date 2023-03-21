package jwt

import (
	"crypto/rsa"
	"fmt"
	"jungle-test/app/internal/domain/entity"
	"jungle-test/app/pkg/apperrors"
	"os"
	"time"

	"github.com/kataras/jwt"
)

type JwtConfigurator struct {
	accessTokenDur  time.Duration
	refreshTokenDur time.Duration
	// sessionMaxAge   time.Duration
	signer        *rsa.PrivateKey
	verifier      *rsa.PublicKey
	signingMethod jwt.Alg
}

func NewJwtConfigurator(accessTokenDur time.Duration, privateKeyPath string) (*JwtConfigurator, error) {
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read private key: %w", err)
	}

	privateKey, err := jwt.ParsePrivateKeyRSA(privateKeyData)
	if err != nil {
		return nil, fmt.Errorf("parse private key rsa: %w", err)
	}

	return &JwtConfigurator{
		accessTokenDur: accessTokenDur,
		signer:         privateKey,
		verifier:       &privateKey.PublicKey,
		signingMethod:  jwt.RS256,
	}, nil
}

func (jc *JwtConfigurator) GenerateAccessToken(user *entity.User) (string, error) {
	now := time.Now()

	accessToken, err := jwt.Sign(jc.signingMethod, jc.signer, AccessToken{

		Claims: jwt.Claims{
			Expiry:   now.Add(jc.accessTokenDur).Unix(),
			IssuedAt: now.Unix(),
		},
		UserClaims: UserClaims{
			UserID:   user.ID,
			Username: user.Username,
		},
	})
	if err != nil {
		return "", apperrors.NewInternal("sign access token", err)
	}

	return string(accessToken), nil
}

func (jc *JwtConfigurator) VerifyAccessToken(token string) (UserClaims, error) {
	verifiedToken, err := jwt.Verify(jc.signingMethod, jc.verifier, []byte(token))
	if err != nil {
		return UserClaims{}, fmt.Errorf("verify token: %w", err)
	}

	var claims UserClaims
	err = verifiedToken.Claims(&claims)
	if err != nil {
		return UserClaims{}, apperrors.NewInternal("parse user claims from access jwt", err)
	}

	return claims, err
}
