package blocks

import (
	"go.mongodb.org/mongo-driver/mongo"

	"bjungle/blockchain-engine/internal/logger"
	"bjungle/blockchain-engine/internal/models"
)

const (
	Postgresql = "postgres"
	MongoDB    = "mongodb"
)

type ServicesBlockRepository interface {
	create(m *Block) error
	update(m *Block) error
	delete(id int64) error
	getByID(id int64) (*Block, error)
	getAll(limit, offSet *int64) ([]*Block, error)
	getHashPrevBlock() (string, error)
	getBlocksById(id int64) (*Block, error)
	existsBlock() bool
}

func FactoryStorage(dbMG *mongo.Database, user *models.User, txID string) ServicesBlockRepository {
	var s ServicesBlockRepository
	engine := MongoDB
	switch engine {
	case MongoDB:
		return newBlockMongodbRepository(dbMG, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
