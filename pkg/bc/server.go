package bc

import (
	"bjungle/blockchain-engine/internal/models"
	"bjungle/blockchain-engine/pkg/bc/blocks"
	"bjungle/blockchain-engine/pkg/bc/blocks_tmp"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	SrvBlocks    blocks.PortsServerBlock
	SrvBlocksTmp blocks_tmp.PortsServerBlockTmp
}

func NewServerBc(db *mongo.Database, user *models.User, txID string) *Server {
	repoBlocks := blocks.FactoryStorage(db, user, txID)
	srvBlocks := blocks.NewBlockService(repoBlocks, user, txID)

	repoBlocksTmp := blocks_tmp.FactoryStorage(db, user, txID)
	srvBlocksTmp := blocks_tmp.NewBlockTmpService(repoBlocksTmp, user, txID)

	return &Server{
		SrvBlocks:    srvBlocks,
		SrvBlocksTmp: srvBlocksTmp,
	}
}
