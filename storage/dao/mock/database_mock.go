package mock

import (
	"github.com/adeo/go-api-skeleton/storage/dao"
	"github.com/stretchr/testify/mock"
)

type DatabaseMock struct {
	mock.Mock
}

func NewDatabaseMock() dao.Database {
	return &DatabaseMock{}
}
