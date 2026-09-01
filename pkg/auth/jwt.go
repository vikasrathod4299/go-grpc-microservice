package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenTTL = time.Hour
	TokenIssuer    = "microservice-auth"
	TokenAudience  = "microservice-gateway"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

var (
	ErrInvalidToken  = errors.New("invalid or expired authentication token")
	ErrInvalidSecret = errors.New("invalid secret key")
)

func GenerateToken(userID, role, secret string) (string, time.Time, error) {
	if len([]byte(secret)) < 32 {
		return "", time.Time{}, ErrInvalidSecret
	}
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)
	claim := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    TokenIssuer,
			Audience:  jwt.ClaimStrings([]string{TokenAudience}),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expiresAt, err
}

func ValidateToken(tokenStr, secret string) (*Claims, error) {
	if len([]byte(secret)) < 32 {
		return nil, ErrInvalidSecret
	}
	token, err := jwt.ParseWithClaims(
		tokenStr, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{
			jwt.SigningMethodHS256.Alg(),
		}),
		jwt.WithIssuer(TokenIssuer),
		jwt.WithAudience(TokenAudience),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
	)
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	if claims.UserID == "" ||
		claims.Role == "" ||
		claims.Subject == "" ||
		claims.Subject != claims.UserID {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
