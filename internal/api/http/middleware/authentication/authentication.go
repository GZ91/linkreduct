package authenticationmiddleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/GZ91/linkreduct/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const SecretKey = "Secret_Key"

type ctxString string

// Authentication возвращает обработчик HTTP, предоставляющий аутентификацию.
// Извлекает идентификатор пользователя (userID) из куки "Authorization".
// Если куки отсутствует, создает новую куку с авторизационными данными.
// В случае ошибок при извлечении или создании куки, возвращает HTTP-статус 500 Internal Server Error.
// Проверяет валидность полученного идентификатора пользователя.
// В случае ошибок при валидации, возвращает HTTP-статус 500 Internal Server Error и продолжает обработку запроса.
// Если идентификатор не валиден, обновляет куку и продолжает обработку запроса.
// Добавляет идентификатор пользователя в контекст запроса под ключом "userID".
// Затем передает управление следующему обработчику в цепочке.
func Authentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Инициализируем переменную для хранения идентификатора пользователя
		userID := ""

		// Извлекаем куку "Authorization" из запроса
		cookie, err := r.Cookie("Authorization")
		if err != nil && err != http.ErrNoCookie {
			// В случае ошибки (не связанной с отсутствием куки), логируем ошибку
			// и возвращаем HTTP-статус 500 Internal Server Error
			logger.Log.Error("authorization", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Если кука "Authorization" отсутствует
		if err == http.ErrNoCookie {
			// Создаем новую куку "Authorization"
			cook := &http.Cookie{}
			cook.Name = "Authorization"

			// Получаем авторизационные данные и идентификатор пользователя
			var value string
			value, userID, err = getAuthorizationForCookie()
			if err != nil {
				// В случае ошибки при получении авторизационных данных, логируем ошибку
				// и возвращаем HTTP-статус 500 Internal Server Error
				logger.Log.Error("authorization", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// Устанавливаем значение куки
			cook.Value = value
			// Устанавливаем куку в ответе
			http.SetCookie(w, cook)
		} else {
			// Если кука "Authorization" присутствует

			// Проверяем валидность идентификатора пользователя
			var ok bool
			userID, ok, err = validGetAuthentication(cookie.Value)
			if err != nil {
				// В случае ошибки при валидации, логируем ошибку
				// и передаем управление следующему обработчику в цепочке
				logger.Log.Error("authorization", zap.Error(err))
				h.ServeHTTP(w, r)
				return
			}

			// Если идентификатор не валиден
			if !ok {
				// Обновляем куку "Authorization" с новыми авторизационными данными
				cookie.Value, userID, err = getAuthorizationForCookie()
				if err != nil {
					// В случае ошибки при получении новых авторизационных данных,
					// логируем ошибку и передаем управление следующему обработчику в цепочке
					logger.Log.Error("authorization", zap.Error(err))
					h.ServeHTTP(w, r)
					return
				}
			}
		}

		// Добавляем идентификатор пользователя в контекст запроса под ключом "userID"
		var userIDCTX models.CtxString = "userID"
		r = r.WithContext(context.WithValue(r.Context(), userIDCTX, userID))

		// Передаем управление следующему обработчику в цепочке
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
			return []byte(SecretKey), nil
		})
	if err != nil {
		return "", false, err
	}

	if !token.Valid {
		return "", false, nil
	}

	return claims.UserID, true, nil
}
