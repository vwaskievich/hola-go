package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"net/http"
	"time"

	"io/ioutil"

	"github.com/gorilla/mux"

	"github.com/diegovillarino/go/tree/victor_user/models"
	s "github.com/diegovillarino/go/tree/victor_user/services/user.services"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/diegovillarino/go/tree/victor_user/awsgo"
	"github.com/diegovillarino/go/tree/victor_user/database"
	"github.com/diegovillarino/go/tree/victor_user/handlers"
	"github.com/diegovillarino/go/tree/victor_user/secretmanager"
)

func EjecutoLambda(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	var res *events.APIGatewayProxyResponse
	awsgo.InicializoAWS()

	if !ValidoParametros() {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error en las variables de entorno. deben incluir 'SecretName', 'BucketName",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil
	} else {
		fmt.Print("Parametros OK")
	}

	SecretModel, err := secretmanager.GetSecret(os.Getenv("SecretName"))
	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error en la lectura de Secret " + err.Error(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil
	} else {
		fmt.Print("Secret OK")
	}

	path := strings.Replace(request.PathParameters["twitter"], os.Getenv("UrlPrefix"), "", -1)

	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("path"), path)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("method"), request.HTTPMethod)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("user"), SecretModel.Username)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("password"), SecretModel.Password)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("host"), SecretModel.Host)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("database"), SecretModel.Database)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("jwtSign"), SecretModel.JWTSign)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("body"), request.Body)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("bucketName"), os.Getenv("BucketName"))

	// Chequeo Conexión a la BD o Conecto la BD

	database.ConectarBD(awsgo.Ctx)

	respAPI := handlers.Manejadores(awsgo.Ctx, request)

	fmt.Println("Sali de Manejadores")
	if respAPI.CustomResp == nil {
		headersResp := map[string]string{
			"Content-Type": "application/json",
		}
		res = &events.APIGatewayProxyResponse{
			StatusCode: respAPI.Status,
			Body:       string(respAPI.Message),
			Headers:    headersResp,
		}
		fmt.Print("Paso por el lambda no ok")
		return res, nil
	} else {
		fmt.Print("Paso por el lambda ok")
		return respAPI.CustomResp, nil
	}	
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := s.Read()

	if err != nil {
		fmt.Print("Se ha presntado un error en la consulta de usuarios")
	}

	if len(users) == 0 {
		fmt.Print("La consulta no retorno datos")
	} else {
		fmt.Print("La lectura finalizo con exitoooooooooooooooo")
	}
	json.NewEncoder(w).Encode(users)
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my API HelloGO")
}

func createUsers(w http.ResponseWriter, r *http.Request) {

	var user models.User
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Insert a valid user")
	}

	json.Unmarshal(reqBody, &user)

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err1 := s.Create(user)

	if err1 != nil {
		fmt.Print("No pudo crearse el usuario")
	} else {
		fmt.Print("Se creo con exito el usuario")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)

}

func updateUsers(w http.ResponseWriter, r *http.Request) {

	var user models.User
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Update a valid user")
	}

	json.Unmarshal(reqBody, &user)

	err1 := s.Update(user, user.ID.Hex())

	if err1 != nil {
		fmt.Print("No pudo actualizarse el usuario")
	} else {
		fmt.Print("Se actualizo el usuario con exito")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)

}

func deleteUser(w http.ResponseWriter, r *http.Request) {

	var user models.User
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "delete a valid user")
	}

	json.Unmarshal(reqBody, &user)

	err1 := s.Delete(user.ID.Hex())

	if err1 != nil {
		fmt.Print("No pudo eliminarse el usuario")
	} else {
		fmt.Print("Se elimino el usuario con exito")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)

}

func main() {

	os.Setenv("AWS_LAMBDA_RUNTIME_API", "runtime.sa-east-1.amazonaws.com")

	myCustomVar := os.Getenv("_LAMBDA_SERVER_PORT")
	fmt.Println("Valor de _LAMBDA_SERVER_PORT:", myCustomVar)

	myCustomVar1 := os.Getenv("AWS_LAMBDA_RUNTIME_API")
	fmt.Println("Valor de AWS_LAMBDA_RUNTIME_API:", myCustomVar1)

	lambda.Start(EjecutoLambda)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users", createUsers).Methods("POST")
	router.HandleFunc("/users", updateUsers).Methods("PUT")
	router.HandleFunc("/users", deleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", router))
	//fmt.Println("conexion a MongoDB")
}

func ValidoParametros() bool {
	_, traeParametro := os.LookupEnv("SecretName")
	if !traeParametro {
		return traeParametro
	}
	_, traeParametro = os.LookupEnv("BucketName")
	if !traeParametro {
		return traeParametro
	}
	_, traeParametro = os.LookupEnv("UrlPrefix")
	if !traeParametro {
		return traeParametro
	}
	return traeParametro
}
