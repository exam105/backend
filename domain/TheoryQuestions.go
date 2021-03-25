package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type TheoryQuestion struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Questions string             `json:"questions,omitempty" bson:"questions"`
	Marks     string             `json:"marks,omitempty" bson:"marks"`
	Answer    string             `json:"answer,omitempty" bson:"answer"`
	Topics    []QuestionTopics   `json:"topics,omitempty" bson:"topics"`	// Referencing from McqQusetion
	Images    []QuestionImages   `json:"images,omitempty" bson:"images"`	// Referencing from McqQusetion
}
