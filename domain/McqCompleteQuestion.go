package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"context"
)

type MCQModel []struct {
	Subject   	string `json:"subject,omitempty" bson:"subject,omitempty"`
	System    	string `json:"system,omitempty" bson:"system,omitempty"`
	Board     	string `json:"board,omitempty" bson:"board,omitempty"`
	Series    	string `json:"series,omitempty" bson:"series,omitempty"`
	Paper     	string `json:"paper,omitempty" bson:"paper,omitempty"`
	Year      	string `json:"year,omitempty" bson:"year,omitempty"`
	Month     	string `json:"month,omitempty" bson:"month,omitempty"`
	Questions 	string `json:"questions,omitempty" bson:"questions,omitempty"`
	Marks     	string `json:"marks,omitempty" bson:"marks,omitempty"`
	Options   	option `json:"options,omitempty" bson:"options,omitempty"`
	Topics    	topic  `json:"topics,omitempty" bson:"topics,omitempty"`
	Images		image  `json:"images,omitempty" bson:"images,omitempty"`
}

type option []struct {
	Option  string  `json:"option" bson:"option"`
	Correct bool    `json:"correct" bson:"correct"`
}

type topic []struct {
	Topic string  `json:"topic" bson:"topic"`
}

type image []struct {
	Imageurl string  `json:"imageurl" bson:"imageurl"`
}

// This struct will be used to display the list of questions to the operator
type DisplayQuestion struct {
	ID       		primitive.ObjectID 	`json:"id,omitempty" bson:"_id,omitempty"`
	Questions 		string 				`json:"questions,omitempty" bson:"questions,omitempty"`
}

// Question Use Case / Service layer
type QuestionUsecase interface {

	GetMetadataById(requestCtx context.Context, username string, useremail string) ([]MetadataBson, error)
	UpdateMetadataById(requestCtx context.Context, metadata MetadataBson, docID string) (int64, error)
	DeleteMetadataById(requestCtx context.Context, docID string) (int64, error)
	GetMCQsByMetadataID(requestCtx context.Context, docID string) ([]DisplayQuestion, error)
	
	SaveMCQ(requestCtx context.Context, questions *MCQModel, username string, useremail string) (error)	
	AddSingleQuestion(requestCtx context.Context, question *Question, metadataID string) (int64)	
	GetQuestion(requestCtx context.Context, questionID string) (Question, error)
	UpdateQuestion(requestCtx context.Context, updatedQuestion Question, questionID string) (int64, error)
	DeleteQuestion(requestCtx context.Context, metadataID string, questionID string) (int64, error)
	
}

// Question Repository represent the question repository contract
type QuestionRepository interface {

	SaveAllQuestions(ctx context.Context, mcq []interface{}) (int64)
	SaveQuestionMetadata(ctx context.Context, mcqMetaData *MetadataBson)
	GetMetadataById(requestCtx context.Context, username string, useremail string) ([]MetadataBson, error)
	UpdateMetadataById(requestCtx context.Context, metadata MetadataBson, docID string) (int64, error)
	DeleteMetadataById(requestCtx context.Context, docID string) (int64, error)
	GetMCQsByMetadataID(requestCtx context.Context, docID string) ([]DisplayQuestion, error)	

	AddSingleQuestion(ctx context.Context, mcq *Question, metadataID string) (int64)
	GetQuestion(requestCtx context.Context, questionID string) (Question, error)
	UpdateQuestion(requestCtx context.Context, updatedQuestion Question, questionID string) (int64, error)
	DeleteQuestion(requestCtx context.Context, metadataID string, questionID string) (int64, error)

}
