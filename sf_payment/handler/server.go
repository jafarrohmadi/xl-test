package handler

import (
	"github.com/xlsmart-api/sf-payment/usecase"
)

type Server struct {
	UseCase usecase.UseCaseInterface
}

type NewServerOptions struct {
	UseCase usecase.UseCaseInterface
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		UseCase: opts.UseCase,
	}
}
