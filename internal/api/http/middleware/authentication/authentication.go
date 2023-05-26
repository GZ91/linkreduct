package authentication

import (
	"context"
	"github.com/GZ91/linkreduct/internal/app/logger"
	"go.uber.org/zap"
	"net/http"
)

func Authentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := ""
		_, err := r.Cookie("Authorization")
		if err != nil && err != http.ErrNoCookie {
			logger.Log.Error("authorization", zap.Error(err))
			h.ServeHTTP(w, r)
			return
		}
		if err == http.ErrNoCookie {
			cook := &http.Cookie{}
			cook.Name = "Authorization"
			var value string
			value, userID, err = getAuthorizationForCookie()
			if err != nil {
				logger.Log.Error("authorization", zap.Error(err))
				h.ServeHTTP(w, r)
				return
			}
			cook.Value = value
			http.SetCookie(w, cook)
		} else {

		}

		r.WithContext(context.WithValue(r.Context(), "userID", userID))

	})
}

func getAuthorizationForCookie() (string, string, error) {

	return "", "", nil
}

func validGetAuthentication(value string) (string, bool, error) {

	return "", false, nil
}
