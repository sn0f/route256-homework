package loms

import (
	"route256/loms/internal/domain"
	api "route256/loms/pkg/loms/v1"
)

type lomsServer struct {
	api.UnimplementedLomsServiceServer
	model *domain.Model
}

func NewLomsServer(model *domain.Model) *lomsServer {
	return &lomsServer{
		model: model,
	}
}
