package middleware

/*
================================================================================
FILE: internal/gateway/middleware/auth.go
================================================================================

PURPOSE:
HTTP Middleware that validates incoming JWT tokens from the `Authorization: Bearer <token>` header.
Extracts user identity (UserID and Role) and stores it in the `r.Context()` for downstream handlers.

LEARNING GO CONCEPTS:
- HTTP middleware wrapping (`func(next http.Handler) http.Handler`).
- `context.WithValue` and `r.WithContext(ctx)`.

WHAT YOU NEED TO IMPLEMENT HERE:
1. `JWTMiddleware(jwtSecret string) func(http.Handler) http.Handler`
2. Extract token from header (`r.Header.Get("Authorization")`).
3. Call `auth.ValidateToken(tokenStr, secret)`.
4. Store claims in context: `ctx := context.WithValue(r.Context(), "user", claims)`.
5. Call `next.ServeHTTP(w, r.WithContext(ctx))`.
================================================================================
*/

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserContextKey contextKey = "user_claims"

func JWTMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
				return
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				http.Error(w, `{"error":"invalid bearer token format"}`, http.StatusUnauthorized)
				return
			}

			// TODO: Validate token string with pkg/auth
			// claims, err := auth.ValidateToken(tokenParts[1], secret)
			// if err != nil { http.Error(...); return }

			// Mock context assignment for now:
			ctx := context.WithValue(r.Context(), UserContextKey, "user-123")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
