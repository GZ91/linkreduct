package authenticationmiddleware

import (
	"context"
	"fmt"
	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

const SecretKey = "Secret_Key"

func Authentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := ""
		cookie, err := r.Cookie("Authorization")
		if err != nil && err != http.ErrNoCookie {
			logger.Log.Error("authorization", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err == http.ErrNoCookie {
			cook := &http.Cookie{}
			cook.Name = "Authorization"
			var value string
			value, userID, err = getAuthorizationForCookie()
			if err != nil {
				logger.Log.Error("authorization", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			cook.Value = value
			http.SetCookie(w, cook)
		} else {
			var ok bool
			userID, ok, err = validGetAuthentication(cookie.Value)
			if err != nil {
				logger.Log.Error("authorization", zap.Error(err))
				h.ServeHTTP(w, r)
				return
			}
			if !ok {
				cookie.Value, userID, err = getAuthorizationForCookie()
				if err != nil {
					logger.Log.Error("authorization", zap.Error(err))
					h.ServeHTTP(w, r)
					return
				}
			}
		}

		r = r.WithContext(context.WithValue(r.Context(), "userID", userID))
		h.ServeHTTP(w, r)
	})
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func getAuthorizationForCookie() (string, string, error) {
	UserID := uuid.New().String()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: UserID,
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", "", err
	}

	return tokenString, UserID, nil
}

func validGetAuthentication(tokenString string) (string, bool, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				strErr := fmt.Sprintf("unexpected signing method: %v", t.Header["alg"])
				logger.Log.Error(strErr)
				return nil, fmt.Errorf(strErr)
			}
			return []byte(SECRET_KEY), nil
		})
	if err != nil {
		return "", false, err
	}

	if !token.Valid {
		return "", false, nil
	}

	return claims.UserID, true, nil
}
