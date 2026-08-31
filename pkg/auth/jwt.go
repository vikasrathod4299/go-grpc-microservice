package auth

/*
================================================================================
PACKAGE: pkg/auth
================================================================================

PURPOSE:
Provides JWT (JSON Web Token) generation and validation utilities.
Used by the API Gateway to verify rider and driver identity on every incoming HTTP/WebSocket request.

LEARNING GO CONCEPTS:
- Working with third-party libraries (e.g. `golang-jwt/jwt/v5`).
- Custom claims structs.
- Symmetric token signing (`SigningMethodHS256`).

WHAT YOU NEED TO IMPLEMENT HERE:
1. Define a `Claims` struct embedding `jwt.RegisteredClaims`:
   - UserID string
   - Role string ("rider" or "driver")
2. Write `GenerateToken(userID string, role string, secret string) (string, error)`
   - Used for testing or auth service to issue tokens.
3. Write `ValidateToken(tokenStr string, secret string) (*Claims, error)`
   - Parses the token string, verifies signature and expiration, and extracts claims.
================================================================================
*/

import (
	"errors"
	// TODO: Uncomment when writing token logic:
	// "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"` // "rider" or "driver"
	// jwt.RegisteredClaims
}

var ErrInvalidToken = errors.New("invalid or expired authentication token")

// TODO: Implement GenerateToken
func GenerateToken(userID, role, secret string) (string, error) {
	// 1. Create Claims with expiry (e.g., 24 hours)
	// 2. Token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 3. return Token.SignedString([]byte(secret))
	return "", nil
}

// TODO: Implement ValidateToken
func ValidateToken(tokenStr, secret string) (*Claims, error) {
	// 1. Parse token string with jwt.ParseWithClaims
	// 2. Validate token.Valid and return claims
	return nil, ErrInvalidToken
}
