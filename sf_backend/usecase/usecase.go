package usecase

import (
	sf_backend_repo "github.com/xlsmart-api/sf-backend/repository"
)

type UseCase struct {
	Repository       sf_backend_repo.RepositoryInterface
	PartnerUseCase   interface{} // TODO: Use proper interface from partner module
	SFPaymentUseCase interface{} // TODO: Use proper interface from sf_payment module
}

type NewUseCaseOptions struct {
	Repository       sf_backend_repo.RepositoryInterface
	PartnerUseCase   interface{} // Optional: for inter-module communication
	SFPaymentUseCase interface{} // Optional: for inter-module communication
}

func NewUseCase(opts NewUseCaseOptions) *UseCase {
	return &UseCase{
		Repository:       opts.Repository,
		PartnerUseCase:   opts.PartnerUseCase,
		SFPaymentUseCase: opts.SFPaymentUseCase,
	}
}
