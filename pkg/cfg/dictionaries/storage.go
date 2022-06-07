package dictionaries

import (
	"go.mongodb.org/mongo-driver/mongo"

	"bjungle/blockchain-engine/internal/logger"
	"bjungle/blockchain-engine/internal/models"
)

const (
	Postgresql = "postgres"
	MongoDB    = "mongodb"
)

type ServicesDictionariesRepository interface {
	create(m *Dictionary) error
	update(m *Dictionary) error
	delete(id int) error
	getByID(id int) (*Dictionary, error)
	getAll() ([]*Dictionary, error)
}

func FactoryStorage(dbMG *mongo.Database, user *models.User, txID string) ServicesDictionariesRepository {
	var s ServicesDictionariesRepository
	engine := MongoDB

	switch engine {
	case MongoDB:
		return newDictionaryMongodbRepository(dbMG, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
