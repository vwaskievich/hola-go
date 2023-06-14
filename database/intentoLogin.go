package database

import (
	"github.com/vwaskievich/hola-go/tree/victor_aws_lambda/models"
	"golang.org/x/crypto/bcrypt"
)

func IntentoLogin(email string, password string) (models.User, bool) {
	usu, encontrado, _ := ChequeoYaExisteUsuario(email)
	if !encontrado {
		return usu, false
	}

	passwordBytes := []byte(password)
	passwordBD := []byte(usu.Password)
	err := bcrypt.CompareHashAndPassword(passwordBD, passwordBytes)
	if err != nil {
		return usu, false
	}
	return usu, true
}
