package user_services_test

import (
	m "../../models"
	userService "../user.services"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

var userId string

func TestCreate(t *testing.T) {

	t.Log("-----------PRUEBA-----------------")
	oid := primitive.NewObjectID()
	userId = oid.Hex()

	user := m.User{
		ID:        oid,
		Name:      "Jorge asdfasdf",
		Email:     "jorge@gmail.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := userService.Create(user)

	if err != nil {
		t.Error("La prueba de persistencia de datos del usuario a fallado")
		t.Fail()
	} else {
		t.Log("La prueba finalizo con exito!", "userId  = ", userId, " - oid = ", oid)
	}
}

func TestRead(t *testing.T) {

	users, err := userService.Read()

	if err != nil {
		t.Error("Se ha presntado un error en la consulta de usuarios")
	}

	if len(users) == 0 {
		t.Error("La consulta no retorno datos")
	} else {
		t.Log("La lectura finalizo con exito", users[0].Name, " - ", users)
	}
}

func TestUpdate(t *testing.T) {

	user := m.User{
		Name:  "Jorge Lopez",
		Email: "jorge.lopez@gmail.com",
	}

	err := userService.Update(user, userId)

	if err != nil {
		t.Error("Error al tratar de actualizar el usuario")
	} else {
		t.Log("La prueba finalizo correctamente")
	}
}

func TestDelete(t *testing.T) {

	err := userService.Delete(userId)

	if err != nil {
		t.Error("Error al tratar de eliminar el usuario")
		t.Fail()
	} else {
		t.Log("La prueba de eliminacion finalizo con exito!")
	}
}
