package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/exam105-UPD/backend/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type loginRepo struct {
	Conn *mongo.Client
}

// This will create an object that represent the login.Repository interface
func NewLoginRepository(Conn *mongo.Client) domain.LoginRepository {
	return &loginRepo{Conn}
}

func (lgRepo *loginRepo) Authenticate(ctx context.Context, username string, useremail string) {

}
func (db *loginRepo) Save(ctx context.Context, DEO_Model *domain.DataEntryOperatorModel) {

	database := db.Conn.Database("exam105")
	dataEntryOperatorCollection := database.Collection("operator_account")
	insertResult, err := dataEntryOperatorCollection.InsertOne(ctx, DEO_Model)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single metadata document: ", insertResult)

	return
}

func (db *loginRepo) GetAllOperators(ctx context.Context) ([]domain.DataEntryOperatorModel, error) {

	database := db.Conn.Database("exam105")
	dataEntryOperatorCollection := database.Collection("operator_account")

	cursor, err := dataEntryOperatorCollection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var operator []domain.DataEntryOperatorModel
	if err = cursor.All(ctx, &operator); err != nil {
		log.Fatal(err)
	}

	if len(operator) == 0 {
		return nil, errors.New("Data Entry Table is empty ")
	}
	fmt.Println(operator)

	return operator, nil
}
