package messages

import (
	"bjungle/blockchain-engine/internal/logger"
	"bjungle/blockchain-engine/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	Postgresql = "postgres"
	MongoDB    = "mongodb"
)

type ServicesMessagesRepository interface {
	create(m *Messages) error
	update(m *Messages) error
	delete(id int) error
	getByID(id int) (*Messages, error)
	getAll() ([]*Messages, error)
}

func FactoryStorage(dbMG *mongo.Database, user *models.User, txID string) ServicesMessagesRepository {
	var s ServicesMessagesRepository
	engine := MongoDB
	switch engine {
	case MongoDB:
		return newMessagesMongodb(dbMG, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
