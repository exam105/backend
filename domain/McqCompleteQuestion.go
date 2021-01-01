package domain

import "context"

type MCQModel []struct {
	Subject   string `json:"subject,omitempty" bson:"subject,omitempty"`
	System    string `json:"system,omitempty" bson:"system,omitempty"`
	Board     string `json:"board,omitempty" bson:"board,omitempty"`
	Series    string `json:"series,omitempty" bson:"series,omitempty"`
	Paper     string `json:"paper,omitempty" bson:"paper,omitempty"`
	Year      string `json:"year,omitempty" bson:"year,omitempty"`
	Month     string `json:"month,omitempty" bson:"month,omitempty"`
	Questions string `json:"questions,omitempty" bson:"questions,omitempty"`
	Marks     string `json:"marks,omitempty" bson:"marks,omitempty"`
	Options   option `json:"options,omitempty" bson:"options,omitempty"`
	Topic     topic  `json:"topic,omitempty" bson:"topic,omitempty"`
}

type option []struct {
	Option  string  `json:"option" bson:"option"`
	ID      float64 `json:"id" bson:"id"`
	Correct bool    `json:"correct" bson:"correct"`
}

type topic []struct {
	Topic string  `json:"topic" bson:"topic"`
	ID    float64 `json:"id" bson:"id"`
}

// Question Use Case / Service layer
type QuestionUsecase interface {
	Save(requestCtx context.Context, questions *MCQModel, username string, useremail string) error
}

// Question Repository represent the question repository contract
type QuestionRepository interface {
	SaveAllQuestions(ctx context.Context, mcq []interface{}) //mcq []Question)
	SaveQuestionMetadata(ctx context.Context, mcqMetaData *MetadataBson)
}
