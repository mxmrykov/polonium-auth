package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type (
	JWTProcessor struct {
		publicKey       *rsa.PublicKey
		privateKey      *rsa.PrivateKey
		signingMethod   jwt.SigningMethod
		access, refresh time.Duration
		issuer          string
	}

	CustomClaims struct {
		UserID   string `json:"user_id"`
		Deployer string `json:"deployer"`
		Session  string `json:"session"`
		jwt.RegisteredClaims
	}
)

func NewRSAJWTProcessor(
	privateKey *rsa.PrivateKey,
	publicKey *rsa.PublicKey,
	access, refresh time.Duration,
	issuer string,
) *JWTProcessor {
	return &JWTProcessor{
		privateKey:    privateKey,
		publicKey:     publicKey,
		signingMethod: jwt.SigningMethodRS256,
		access:        access,
		refresh:       refresh,
		issuer:        issuer,
	}
}

func (j *JWTProcessor) GenerateAccessToken(user, session string) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		UserID:  user,
		Session: session,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.access)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   "access",
		},
	}

	token := jwt.NewWithClaims(j.signingMethod, claims)
	return token.SignedString(j.privateKey)
}

func (j *JWTProcessor) GenerateRefreshToken(session string, userID string) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		UserID:  userID,
		Session: session,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.refresh)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   "refresh",
		},
	}

	token := jwt.NewWithClaims(j.signingMethod, claims)
	return token.SignedString(j.privateKey)
}

func (j *JWTProcessor) TokenVerify(tokenString string) (*CustomClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token is empty")
	}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != j.signingMethod.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		return j.publicKey, nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	if j.issuer != "" && claims.Issuer != j.issuer {
		return nil, fmt.Errorf("invalid token issuer: %s", claims.Issuer)
	}

	return claims, nil
}
