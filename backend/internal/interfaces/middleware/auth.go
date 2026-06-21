package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"lokalsicht/internal/domain/account"
)

type contextKey string

const UserKey contextKey = "user"

func RequireAuth(secret string, userRepo account.AccountRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				writeError(w, http.StatusUnauthorized, "missing token")
				return
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				writeError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				writeError(w, http.StatusUnauthorized, "invalid claims")
				return
			}
			email, _ := claims["email"].(string)
			if email == "" {
				writeError(w, http.StatusUnauthorized, "no email in token")
				return
			}

			user, err := userRepo.FindByEmail(r.Context(), email)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "internal error")
				return
			}
			if user == nil {
				// Create user and account on first request
				name, _ := claims["name"].(string)
				acc := &account.Account{Name: name, Plan: account.PlanBasic}
				if createErr := userRepo.Create(r.Context(), acc); createErr != nil {
					writeError(w, http.StatusInternalServerError, "failed to create account")
					return
				}
				user = &account.User{Email: email, Name: name, AccountID: acc.ID, Role: "owner"}
				if createErr := userRepo.CreateUser(r.Context(), user); createErr != nil {
					writeError(w, http.StatusInternalServerError, "failed to create user")
					return
				}
				// Reload to get the Account preloaded
				user, _ = userRepo.FindByEmail(r.Context(), email)
			}

			ctx := context.WithValue(r.Context(), UserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUser(ctx context.Context) *account.User {
	user, _ := ctx.Value(UserKey).(*account.User)
	return user
}

func writeError(w http.ResponseWriter, status int, detail string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"error":"%s"}`, detail)
}
