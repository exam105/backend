package domain

import "context"

type UserLogin struct {
	Username string
	Email    string
}

// Login Use Case / Service layer
type LoginUsecase interface {
	Authenticate(ctx context.Context)
	// Save(ctx context.Context) error
}

// ArticleRepository represent the article's repository contract
type LoginRepository interface {
	Authenticate(ctx context.Context, username string, useremail string)
}
