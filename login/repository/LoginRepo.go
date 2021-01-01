package repository

import (
	"context"

	"github.com/exam105-UPD/backend/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type loginRepo struct {
	Conn *mongo.Client
}

// This will create an object that represent the login.Repository interface
func NewLoginRepository(Conn *mongo.Client) domain.LoginRepository {
	return &loginRepo{Conn}
}

func (lgRepo *loginRepo) Authenticate(ctx context.Context, username string, useremail string) {

}
