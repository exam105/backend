package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/exam105-UPD/backend/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (db *loginRepo) Update(ctx context.Context, dataEntryOperator *domain.DataEntryOperatorModel, objectId primitive.ObjectID) (int64, error) {

	database := db.Conn.Database("exam105")
	deoCollection := database.Collection("operator_account")

	// Use it's ID to replace
	filter := bson.M{"_id": objectId}
	// Create a replacement object using the existing object
	replacementObj := dataEntryOperator
	///replacementObj.Id = objectId
	replacementObj.Username = dataEntryOperator.Username
	replacementObj.Email = dataEntryOperator.Email
	updateResult, err := deoCollection.ReplaceOne(ctx, filter, replacementObj)

	if err != nil {
		log.Println(err)
	}

	fmt.Printf(
		"Match count: %d, updated: %d, deleted: %v",
		updateResult.MatchedCount,
		updateResult.ModifiedCount,
		updateResult.UpsertedID,
		//res.DeletedCount,
	)

	return updateResult.ModifiedCount, nil
}

func (db *loginRepo) Delete(ctx context.Context, objectId primitive.ObjectID) (int64, error) {

	database := db.Conn.Database("exam105")
	deoCollection := database.Collection("operator_account")

	fmt.Println("DeleteOne TYPE:", reflect.TypeOf(deoCollection))

	// Use it's ID to replace
	filter := bson.M{"_id": objectId}
	deleteResult, err := deoCollection.DeleteOne(ctx, filter)

	if err != nil {
		log.Println(err)
	}

	fmt.Printf("deleted: %v", deleteResult.DeletedCount)

	return deleteResult.DeletedCount, nil
}
