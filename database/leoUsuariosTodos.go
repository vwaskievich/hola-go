package database

import (
	"context"
	"fmt"
	"time"

	"github.com/diegovillarino/go/tree/victor_user/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
LeoUsuariosTodos Lee los usuarios registrados en el sistema, si se recibe "R" en quienes

	trae solo los que se relacionan conmigo
*/
func LeoUsuariosTodos(ID string, page int64, search string, tipo string) ([]*models.User, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := MongoCN.Database(DatabaseName)
	col := db.Collection("users")

	var results []*models.User

	findOptions := options.Find()
	findOptions.SetLimit(20)
	findOptions.SetSkip((page - 1) * 20)

	query := bson.M{
		"nombre": bson.M{"$regex": `(?i)` + search},
	}

	cur, err := col.Find(ctx, query, findOptions)
	if err != nil {
		return results, false
	}

	for cur.Next(ctx) {
		var s models.User
		err := cur.Decode(&s)
		if err != nil {
			fmt.Println("Decode = " + err.Error())
			return results, false
		}
	}

	err = cur.Err()
	if err != nil {
		fmt.Println("cur.Err() = " + err.Error())
		return results, false
	}
	cur.Close(ctx)
	return results, true
}