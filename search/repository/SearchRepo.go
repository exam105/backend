package repository

import (
	"log"
	"github.com/exam105-UPD/backend/logging"
	"go.mongodb.org/mongo-driver/bson"
	"fmt"
	"github.com/exam105-UPD/backend/domain"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type searchRepo struct {
	Conn *mongo.Client
}

// This will create an object that represent the search.Repository interface
func NewSearchRepository(Conn *mongo.Client) domain.SearchRepository {
	return &searchRepo{Conn}
}

func (db *searchRepo) SearchByDate(ctx context.Context, searchCriteria *domain.SearchParameterByDate) ([]domain.SearchResult_Paper, error){

	database := db.Conn.Database("exam105")
	metadataCollection := database.Collection("metadata")
	cursor, err := metadataCollection.Find(
		ctx,
		bson.D{
			{"subject", searchCriteria.Subject},			
			{"system", searchCriteria.System},
			{"date", searchCriteria.Date},
		})
	if err != nil {
		log.Println( logging.MSG_SearchFailed, err.Error())
		return []domain.SearchResult_Paper{}, fmt.Errorf(logging.MSG_SearchFailed, nil)
	}

	defer cursor.Close(ctx)
	var paperList []domain.SearchResult_Paper
	for cursor.Next(ctx) {
		var paper domain.SearchResult_Paper
		cursor.Decode(&paper)
		paperList = append(paperList, paper)
	}
	
	// fmt.Printf("Search Result -->> \n %s", paperList)
	
	return paperList, nil	
}

func (db *searchRepo) SearchByDateRange(ctx context.Context, searchCriteria *domain.SearchParameterByDateRange) ([]domain.SearchResult_Paper, error){

	fromDate := searchCriteria.FromDate
	toDate := searchCriteria.ToDate

	database := db.Conn.Database("exam105")
	metadataCollection := database.Collection("metadata")
	cursor, err := metadataCollection.Find(
		ctx,
		bson.D{
			{"subject", searchCriteria.Subject},
			{"system", searchCriteria.System},
			{"date", bson.D{{"$gte", fromDate}, {"$lt", toDate}} },
		})
	
	if err != nil {
		log.Println( logging.MSG_SearchFailed, err)
		return []domain.SearchResult_Paper{}, fmt.Errorf(logging.MSG_SearchFailed + err.Error(), nil)
	}

	defer cursor.Close(ctx)
	var paperList []domain.SearchResult_Paper
	for cursor.Next(ctx) {
		var paper domain.SearchResult_Paper
		cursor.Decode(&paper)
		paperList = append(paperList, paper)
	}
	
	fmt.Printf("Search Result -->> \n %s", paperList)
	
	return paperList, nil	
}