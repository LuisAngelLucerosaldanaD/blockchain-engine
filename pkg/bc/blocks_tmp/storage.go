package blocks_tmp

import (
	"bjungle/blockchain-engine/internal/logger"
	"bjungle/blockchain-engine/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	MongoDB = "mongodb"
)

type ServicesBlockTmpRepository interface {
	create(m *BlockTmp) error
	update(m *BlockTmp) error
	delete(id int64) error
	getByID(id int64) (*BlockTmp, error)
	getAll() ([]*BlockTmp, error)
	getBlockUnCommit() (*BlockTmp, error)
	getBlockTwoCommit() (*BlockTmp, error)
	GetCountTransactionByID(block int64) int64
}

func FactoryStorage(dbMG *mongo.Database, user *models.User, txID string) ServicesBlockTmpRepository {
	var s ServicesBlockTmpRepository
	engine := MongoDB
	switch engine {
	case MongoDB:
		return newBlockTempMongodbRepository(dbMG, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
