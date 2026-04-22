package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenUtil struct {
	Token string
}

func NewTokenUtil(Secret string) TokenUtil {
	return TokenUtil{
		Token: Secret,
	}
}

func (t *TokenUtil) NewJwtToken(UserID string) (string, error) {

	claims := jwt.MapClaims{
		"ID":  UserID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	}

	Token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	JwtToken, err := Token.SignedString([]byte(t.Token))

	if err != nil {
		return " ", err
	}

	return JwtToken, nil

}

func (t *TokenUtil) ParsedToken(TokenString string) (string, error) {

	TokenParse, err := jwt.Parse(TokenString, func(tokenParam *jwt.Token) (interface{}, error) {

		if _, ok := tokenParam.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode tanda tangan tidak valid/mencurigakan: %v", tokenParam.Header["alg"])
		}

		return []byte(t.Token), nil

	})

	if err != nil {
		return " ", fmt.Errorf("gagal memproses toke %v", err)
	}

	claims, ok := TokenParse.Claims.(jwt.MapClaims)

	if ok && TokenParse.Valid {
		id, ok := claims["ID"].(string)

		if !ok {
			return " ", fmt.Errorf("User tidak ditemukan atau format salah %v", err)
		}
		return id, nil
	}

	return " ", err

}
