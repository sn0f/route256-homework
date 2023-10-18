package services

import (
	"route256/checkout/internal/domain/services/mocks"
	"testing"

	"github.com/gojuno/minimock/v3"
)

type mocksService struct {
	repo        *mocks.CartRepositoryMock
	products    *mocks.ProductManagerMock
	orders      *mocks.OrderManagerMock
	txManager   *mocks.TransactionManagerMock
	workerCount int
}

func newMocksService(t *testing.T, workerCount int) *mocksService {
	mc := minimock.NewController(t)
	return &mocksService{
		repo:        mocks.NewCartRepositoryMock(mc),
		products:    mocks.NewProductManagerMock(mc),
		orders:      mocks.NewOrderManagerMock(mc),
		txManager:   mocks.NewTransactionManagerMock(mc),
		workerCount: workerCount,
	}
}
