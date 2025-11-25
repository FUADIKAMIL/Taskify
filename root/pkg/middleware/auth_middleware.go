package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/yourname/taskify/pkg/auth"
)

type ctxKey string

const UserIDKey ctxKey = "user_id"

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        h := r.Header.Get("Authorization")
        if h == "" {
            http.Error(w, "missing authorization", http.StatusUnauthorized)
            return
        }
        parts := strings.SplitN(h, " ", 2)
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            http.Error(w, "invalid authorization header", http.StatusUnauthorized)
            return
        }
        token := parts[1]
        claims, err := auth.ParseToken(token)
        if err != nil {
            http.Error(w, "invalid token", http.StatusUnauthorized)
            return
        }
        ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
