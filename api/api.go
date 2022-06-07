package api

import (
	"bjungle/blockchain-engine/internal/db"
	"github.com/google/uuid"
)

func Start(port int) {
	dbMg := db.GetMgConnection()

	server := newServer(port, dbMg, uuid.New().String())
	server.Start()
}
