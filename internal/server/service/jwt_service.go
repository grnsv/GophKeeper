package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/grnsv/GophKeeper/internal/server/interfaces"
)

type JWTService struct {
	secret        []byte
	signingMethod jwt.SigningMethod
}

func NewJWTService(secret string) interfaces.JWTService {
	return &JWTService{
		secret:        []byte(secret),
		signingMethod: jwt.SigningMethodHS256,
	}
}

func (s *JWTService) BuildJWT(userID string) (string, error) {
	now := jwt.NewNumericDate(time.Now())
	token := jwt.NewWithClaims(s.signingMethod, jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		NotBefore: now,
		IssuedAt:  now,
	})

	return token.SignedString(s.secret)
}

func (s *JWTService) ParseJWT(token string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if t.Method == nil || t.Method.Alg() != s.signingMethod.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return "", err
	}
	if !jwtToken.Valid {
		return "", fmt.Errorf("token is not valid: %v", token)
	}

	return claims.Subject, nil
}
