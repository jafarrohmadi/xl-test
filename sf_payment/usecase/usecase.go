package usecase

import (
	sf_payment_repo "github.com/xlsmart-api/sf-payment/repository"
)

type UseCase struct {
	Repository sf_payment_repo.RepositoryInterface
}

type NewUseCaseOptions struct {
	Repository sf_payment_repo.RepositoryInterface
}

func NewUseCase(opts NewUseCaseOptions) *UseCase {
	return &UseCase{
		Repository: opts.Repository,
	}
}
