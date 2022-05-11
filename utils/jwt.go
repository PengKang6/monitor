package utils

import (
	"github.com/dgrijalva/jwt-go"
	"monitor/config"
	"time"
)

type Claims struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,required"`
	jwt.StandardClaims
}

var (
	JwtSecret  []byte
	Issuer     string
	ExpireTime int
)

func JWTGenerate(id string, name string) (string, error) {
	now := time.Now()
	expireTime := now.Add(time.Second * time.Duration(ExpireTime))
	issuer := "pk"
	claims := Claims{
		ID:   id,
		Name: name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    issuer,
		},
	}

	token, e := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(JwtSecret)
	return token, e
}

func JWTParse(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (i interface{}, e error) {
		return JwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

func InitJWT(config config.JWTConfig) {
	JwtSecret = []byte(config.JWTSecret)
	Issuer = config.Issuer
	ExpireTime = config.ExpireTime
}
