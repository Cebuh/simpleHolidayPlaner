package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cebuh/simpleHolidayPlaner/config"
	"github.com/cebuh/simpleHolidayPlaner/types"
	"github.com/cebuh/simpleHolidayPlaner/utils"
	"github.com/golang-jwt/jwt"
)

const UserKey string = "userId"

func CreateJWT(secret []byte, userId string) (string, error) {

	expiration := time.Second * time.Duration(config.Envs.JWTExpireTimeInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    userId,
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func Require(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request)  {
		tokenString := extractTokenFromRequest(r)
		token, err := validateToken(tokenString)

		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userId := claims["userID"].(string)
		u, err := store.GetUserById(userId)
		if err != nil {
			log.Printf("failed to get user by id: %v", err)
			permissionDenied(w)
			return
		}

		ctx := r.Context()
		ctx  = context.WithValue(ctx, UserKey, u.Id)
		r = r.WithContext(ctx)

		handlerFunc(w,r)
	}
}

func permissionDenied(w http.ResponseWriter){
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func extractTokenFromRequest(r *http.Request) (string) {
	token := r.Header.Get("Authorization")
	if token != "" {
		return token
	}

	return ""
}

func validateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret),nil
	})
}

func GetUserIdFromContext(ctx context.Context) string {
	userId, ok := ctx.Value(UserKey).(string)
	if !ok {
		return ""
	}

	return userId
}
