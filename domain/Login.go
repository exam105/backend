package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DataEntryOperatorModel struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Username string             `json:"username,omitempty" bson:"username"`
	Email    string             `json:"email,omitempty" bson:"email"`
}

// Login Use Case / Service layer
type LoginUsecase interface {
	Authenticate(ctx context.Context, username string, useremail string) error
	Save(ctx context.Context, dataEntryOperator *DataEntryOperatorModel) error
	GetAllOperators(ctx context.Context) ([]DataEntryOperatorModel, error)
	Update(ctx context.Context, dataEntryOperator *DataEntryOperatorModel, objectId primitive.ObjectID) (int64, error)
	Delete(ctx context.Context, objectId primitive.ObjectID) (int64, error)
}

// ArticleRepository represent the article's repository contract
type LoginRepository interface {
	Authenticate(ctx context.Context, username string, useremail string) error
	Save(ctx context.Context, dataEntryOperator *DataEntryOperatorModel)
	GetAllOperators(ctx context.Context) ([]DataEntryOperatorModel, error)
	Update(ctx context.Context, dataEntryOperator *DataEntryOperatorModel, objectId primitive.ObjectID) (int64, error)
	Delete(ctx context.Context, objectId primitive.ObjectID) (int64, error)
}
