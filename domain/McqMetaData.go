package domain

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MetadataBson struct {
	ID             	primitive.ObjectID 	`json:"id,omitempty" bson:"_id,omitempty"`
	Subject        	string             	`json:"subject,omitempty" bson:"subject,omitempty"`
	System         	string             	`json:"system,omitempty" bson:"system,omitempty"`
	Board          	string             	`json:"board,omitempty" bson:"board,omitempty"`
	Series         	string             	`json:"series,omitempty" bson:"series,omitempty"`
	Paper          	string             	`json:"paper,omitempty" bson:"paper,omitempty"`
	Date           	time.Time          	`json:"date,omitempty" bson:"date"`
	Username       	string             	`json:"username,omitempty" bson:"username"`
	Useremail      	string             	`json:"useremail,omitempty" bson:"useremail"`
	QuestionHexIds 	[]string           	`json:"question_hex_ids,omitempty" bson:"question_hex_ids"`
	IsTheory       	bool               	`json:"is_theory,omitempty" bson:"is_theory"`
	Reference   	string    		  	`json:"reference,omitempty" bson:"reference"`
}

type Metadata struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Subject        string             `json:"subject,omitempty" bson:"subject,omitempty"`
	System         string             `json:"system,omitempty" bson:"system,omitempty"`
	Board          string             `json:"board,omitempty" bson:"board,omitempty"`
	Series         string             `json:"series,omitempty" bson:"series,omitempty"`
	Paper          string             `json:"paper,omitempty" bson:"paper,omitempty"`
	Date           time.Time          `json:"date,omitempty" bson:"date"`
	QuestionHexIds []string           `json:"question_hex_ids,omitempty" bson:"question_hex_ids"`
	IsTheory       bool               `json:"is_theory,omitempty" bson:"is_theory"`
	Reference      string    		  `json:"reference,omitempty" bson:"reference"`	
}