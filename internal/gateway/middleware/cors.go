package middleware

/*
================================================================================
FILE: internal/gateway/middleware/cors.go
================================================================================

PURPOSE:
Enables Cross-Origin Resource Sharing (CORS) so web & mobile clients can communicate with the API Gateway.
================================================================================
*/

import "net/http"

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// CORS is an alias for CORSMiddleware so you can use customMiddleware.CORS in chi router
func CORS(next http.Handler) http.Handler {
	return CORSMiddleware(next)
}
