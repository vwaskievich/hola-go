package jwt

import (
	"context"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/vwaskievich/hola-go/tree/victor_aws_lambda/models"
)

func GeneroJWT(ctx context.Context, t models.User) (string, error) {

	jwtSign := ctx.Value(models.Key("jwtSign")).(string)
	miClave := []byte(jwtSign)

	payload := jwt.MapClaims{
		"email": t.Email,
		"name":  t.Name,
		"_id":   t.ID.Hex(),
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenStr, err := token.SignedString(miClave)
	if err != nil {
		return tokenStr, err
	}
	return tokenStr, nil
}
