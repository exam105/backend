package usecase

import (
	"context"
	"time"

	"github.com/exam105-UPD/backend/domain"
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

func (login *loginUsecase) Authenticate(ctx context.Context) {}
