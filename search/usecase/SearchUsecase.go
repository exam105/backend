package usecase

import (
	"context"
	"time"
	"github.com/exam105-UPD/backend/domain"
)

type searchUsecase struct {
	searchRepo   domain.SearchRepository
	contextTimeout time.Duration
}

// NewSearchUsecase will create new an searchUsecase object representation of domain.SearchUsecase interface in Search.go file
func NewSearchUsecase(srRepo domain.SearchRepository, timeout time.Duration) domain.SearchUsecase {
	return &searchUsecase{
		searchRepo:   srRepo,
		contextTimeout: timeout,
	}
}

func (searchUC *searchUsecase) SearchByDate(requestCtx context.Context, searchCriteria *domain.SearchParameterByDate) ([]domain.SearchResult_Paper, error){

	ctx, cancel := context.WithTimeout(requestCtx, searchUC.contextTimeout)
	defer cancel()

	return searchUC.searchRepo.SearchByDate(ctx, searchCriteria) 
}
