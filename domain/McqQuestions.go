package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Question struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Questions string             `json:"questions,omitempty" bson:"questions"`
	Marks     string             `json:"marks,omitempty" bson:"marks"`
	Options   []QuestionOptions  `json:"options,omitempty" bson:"options"`
	Topics    []QuestionTopic    `json:"topics,omitempty" bson:"topics"`
}

type QuestionOptions struct {
	Option  string `json:"option" bson:"option"`
	Correct bool   `json:"correct" bson:"correct"`
}

type QuestionTopic struct {
	Topic string `json:"topic" bson:"topic"`
}
