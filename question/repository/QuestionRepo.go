package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/exam105-UPD/backend/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type questionRepo struct {
	Conn *mongo.Client
}

// This will create an object that represent the question.Repository interface
func NewQuestionRepository(Conn *mongo.Client) domain.QuestionRepository {
	return &questionRepo{Conn}
}

func (db *questionRepo) SaveQuestionMetadata(ctx context.Context, qsMetaData *domain.MetadataBson) {

	database := db.Conn.Database("exam105")
	metadataCollection := database.Collection("metadata")
	insertResult, err := metadataCollection.InsertOne(ctx, qsMetaData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single metadata document: ", insertResult)

	return
}

func (db *questionRepo) SaveAllQuestions(ctx context.Context, questions []interface{}) {

	database := db.Conn.Database("exam105")
	questionsCollection := database.Collection("questions")

	// create the slice of write models
	var writes []mongo.WriteModel
	for _, ins := range questions {
		model := mongo.NewInsertOneModel().SetDocument(ins)
		writes = append(writes, model)
	}

	// run bulk write
	res, err := questionsCollection.BulkWrite(ctx, writes)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"insert: %d, updated: %d, deleted: %d",
		res.InsertedCount,
		res.ModifiedCount,
		res.DeletedCount,
	)

	//insertResult, err := questionsCollection.InsertMany(ctx, questions)
	if err != nil {
		log.Fatal(err)
	}
	return
}
