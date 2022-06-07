package cfg

import (
	"bjungle/blockchain-engine/internal/models"
	dictionaries2 "bjungle/blockchain-engine/pkg/cfg/dictionaries"
	messages2 "bjungle/blockchain-engine/pkg/cfg/messages"

	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	SrvDictionaries dictionaries2.PortsServerDictionaries
	SrvMessage      messages2.PortsServerMessages
}

func NewServerCfg(dbMg *mongo.Database, user *models.User, txID string) *Server {

	repoDictionaries := dictionaries2.FactoryStorage(dbMg, user, txID)
	srvDictionaries := dictionaries2.NewDictionariesService(repoDictionaries, user, txID)

	repoMessage := messages2.FactoryStorage(dbMg, user, txID)
	srvMessage := messages2.NewMessagesService(repoMessage, user, txID)

	return &Server{
		SrvDictionaries: srvDictionaries,
		SrvMessage:      srvMessage,
	}
}
