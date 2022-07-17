package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"io/ioutil"
	"time"
)

func VerifyToken(token string) (bool, error) {
	pubAccess, err := ioutil.ReadFile("public.pem")

	if err != nil {
		return false, err
	}

	accessPubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubAccess)
	if err != nil {
		return false, err
	}

	parseToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return accessPubKey, nil
	})
	if err != nil {
		return false, err
	}
	claims, ok := parseToken.Claims.(jwt.MapClaims)
	if ok && parseToken.Valid {
		atExp, ok := claims["exp"].(float64)
		if !ok {
			return false, fmt.Errorf("error while getting time of expiring")
		}
		atExpires := int64(atExp)
		if atExpires < time.Now().Unix() {
			return false, fmt.Errorf("token is expired")
		} else {
			return true, nil
		}
	} else {
		return false, fmt.Errorf("token is not valid")
	}
}
