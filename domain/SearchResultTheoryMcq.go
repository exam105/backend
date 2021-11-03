package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

//This can contain Theory and MCQ questions and it is used in search result.

type SearchResult_TheoryMcq struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Question  string             `json:"question,omitempty" bson:"question"`
	Marks     string             `json:"marks,omitempty" bson:"marks"`
	Answer    string             `json:"answer,omitempty" bson:"answer"`
	Options   []QuestionOptions  `json:"options,omitempty" bson:"options"`  // Referencing from McqQusetion
	Topics    []QuestionTopics   `json:"topics,omitempty" bson:"topics"`	// Referencing from McqQusetion
	Images    []QuestionImages   `json:"images,omitempty" bson:"images"`	// Referencing from McqQusetion
	
}