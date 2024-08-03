package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/reymooy27/arena-backend/auth-service/handlers"
	"github.com/reymooy27/arena-backend/auth-service/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				utils.JSONResponse(w, http.StatusUnauthorized, "Unauthorized")
				return
			}
			utils.JSONResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		tokenString := cookie.Value
		claims := &handlers.Claim{}

		secret := os.Getenv("SECRET")
		if secret == "" {
			log.Println("SECRET not set")
		}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil {
			if err != jwt.ErrSignatureInvalid {
				utils.JSONResponse(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			utils.JSONResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		if !token.Valid {
			utils.JSONResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "user", claims))

		next.ServeHTTP(w, r)
	})
}
