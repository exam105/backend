package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/exam105-UPD/backend/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type loginUsecase struct {
	loginRepo      domain.LoginRepository
	contextTimeout time.Duration
}

// NewLoginUsecase will create new an loginUsecase object representation of domain.LoginUsecase interface
func NewLoginUsecase(logRepo domain.LoginRepository, timeout time.Duration) domain.LoginUsecase {
	return &loginUsecase{
		loginRepo:      logRepo,
		contextTimeout: timeout,
	}
}

func (loginUC *loginUsecase) Authenticate(ctx context.Context) {}

func (loginUC *loginUsecase) Save(ctx context.Context, DEO_Model *domain.DataEntryOperatorModel) (err error) {

	ctx, cancel := context.WithTimeout(ctx, loginUC.contextTimeout)
	defer cancel()

	//userLoginModel := new(domain.UserLoginModel)
	DEO_Model.Id = primitive.NewObjectID()
	fmt.Println(DEO_Model)
	loginUC.loginRepo.Save(ctx, DEO_Model)
	return
}

func (loginUC *loginUsecase) GetAllOperators(ctx context.Context) ([]domain.DataEntryOperatorModel, error) {

	ctx, cancel := context.WithTimeout(ctx, loginUC.contextTimeout)
	defer cancel()

	return loginUC.loginRepo.GetAllOperators(ctx)
}
