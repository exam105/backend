package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Question struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Questions string             `json:"questions,omitempty" bson:"questions"`
	Marks     string             `json:"marks,omitempty" bson:"marks"`
	Options   []QuestionOptions  `json:"options,omitempty" bson:"options"`
	Topics    []QuestionTopics    `json:"topics,omitempty" bson:"topics"`
	Images    []QuestionImages    `json:"images,omitempty" bson:"images"`
}

type QuestionOptions struct {
	Option  string `json:"option" bson:"option"`
	Correct bool   `json:"correct" bson:"correct"`
}

type QuestionTopics struct {
	Topic string `json:"topic" bson:"topic"`
}

type QuestionImages struct {
	Imageurl string `json:"imageurl" bson:"imageurl"`
}
