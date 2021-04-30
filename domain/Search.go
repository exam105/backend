package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)


// This struct will be populated with the search parameter sent by the frontend
type SearchParameterByDate struct {
	Subject string    `json:"subject,omitempty" bson:"subject,omitempty"`
	System  string    `json:"system,omitempty" bson:"system,omitempty"`
	Board   string    `json:"board,omitempty" bson:"board,omitempty"`
	Date    time.Time `json:"date,omitempty" bson:"date"`
}

// This struct is a search result which backend will send
type SearchResult_Paper struct {
	ID       		primitive.ObjectID 	`json:"id,omitempty" bson:"_id,omitempty"`
	Question 		string 				`json:"question,omitempty" bson:"question,omitempty"`
	Subject 		string    			`json:"subject,omitempty" bson:"subject,omitempty"`
	System  		string    			`json:"system,omitempty" bson:"system,omitempty"`
	Board   		string    			`json:"board,omitempty" bson:"board,omitempty"`
	Date    		time.Time 			`json:"date,omitempty" bson:"date"`

}

// Question Use Case / Service layer
type SearchUsecase interface {

	SearchByDate(requestCtx context.Context, searchCriteria *SearchParameterByDate) ([]SearchResult_Paper, error)
	
}

// Search Repository represent the search repository contract
type SearchRepository interface {

	SearchByDate(ctx context.Context, searchCriteria *SearchParameterByDate) ([]SearchResult_Paper, error)

}