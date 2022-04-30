package usecase

import (
	"context"
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

func (loginUC *loginUsecase) Authenticate(ctx context.Context, username string, useremail string) (err error) {

	ctx, cancel := context.WithTimeout(ctx, loginUC.contextTimeout)
	defer cancel()

	return loginUC.loginRepo.Authenticate(ctx, username, useremail)

}

func (loginUC *loginUsecase) Save(ctx context.Context, DEO_Model *domain.DataEntryOperatorModel) (err error) {

	ctx, cancel := context.WithTimeout(ctx, loginUC.contextTimeout)
	defer cancel()

	//userLoginModel := new(domain.UserLoginModel)
	DEO_Model.ID = primitive.NewObjectID()
	//fmt.Println(DEO_Model)
	loginUC.loginRepo.Save(ctx, DEO_Model)
	return
}

func (loginUC *loginUsecase) GetAllOperators(ctx context.Context) ([]domain.DataEntryOperatorModel, error) {

	ctx, cancel := context.WithTimeout(ctx, loginUC.contextTimeout)
	defer cancel()

	return loginUC.loginRepo.GetAllOperators(ctx)
}

func (loginUC *loginUsecase) Update(ctx context.Context, dataEntryOperator *domain.DataEntryOperatorModel, objID primitive.ObjectID) (int64, error) {

	ctx, cancel := context.WithTimeout(ctx, loginUC.contextTimeout)
	defer cancel()

	return loginUC.loginRepo.Update(ctx, dataEntryOperator, objID)

}

func (loginUC *loginUsecase) Delete(ctx context.Context, objID primitive.ObjectID) (int64, error) {

	ctx, cancel := context.WithTimeout(ctx, loginUC.contextTimeout)
	defer cancel()

	return loginUC.loginRepo.Delete(ctx, objID)

}
