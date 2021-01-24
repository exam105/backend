package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DataEntryOperatorModel struct {
	Id       primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Username string             `json:"username,omitempty" bson:"username"`
	Email    string             `json:"email,omitempty" bson:"email"`
}

// Login Use Case / Service layer
type LoginUsecase interface {
	Authenticate(ctx context.Context)
	Save(ctx context.Context, dataEntryOperator *DataEntryOperatorModel) error
	GetAllOperators(ctx context.Context) ([]DataEntryOperatorModel, error)
}

// ArticleRepository represent the article's repository contract
type LoginRepository interface {
	Authenticate(ctx context.Context, username string, useremail string)
	Save(ctx context.Context, dataEntryOperator *DataEntryOperatorModel)
	GetAllOperators(ctx context.Context) ([]DataEntryOperatorModel, error)
}
