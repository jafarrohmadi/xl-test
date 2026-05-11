package usecase

import (
	partner_repo "github.com/xlsmart-api/partner/repository"
)

type UseCase struct {
	Repository partner_repo.RepositoryInterface
}

type NewUseCaseOptions struct {
	Repository partner_repo.RepositoryInterface
}

func NewUseCase(opts NewUseCaseOptions) *UseCase {
	return &UseCase{
		Repository: opts.Repository,
	}
}
