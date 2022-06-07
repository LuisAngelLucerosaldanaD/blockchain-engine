package db

import (
	"bjungle/blockchain-engine/internal/env"
	"bjungle/blockchain-engine/internal/logger"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"sync"
)

var (
	once sync.Once
	mg   *mongo.Database
)

func init() {
	once.Do(func() {
		setConnection()
	})
}

func setConnection() {
	c := env.NewConfiguration()
	if c.DB.Engine == "mongodb" {
		clientOptions := options.Client().ApplyURI(connectionString())
		connection, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			logger.Error.Printf("no se puede conectar a la base de datos MongoDB: %v", err)
			panic(err)
		}

		err = connection.Ping(context.TODO(), nil)
		if err != nil {
			logger.Error.Printf("no se cumple prueba de conexión a la base de datos MongoDB: %v", err)
			panic(err)
		}
		mg = connection.Database(c.DB.Name)
		return
	}
}

func connectionString() string {
	var host, database, username, password string
	c := env.NewConfiguration()
	host = c.DB.Server
	database = c.DB.Name
	username = c.DB.User
	password = c.DB.Password

	switch strings.ToLower(c.DB.Engine) {
	case "mongodb":
		return fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority", username, password, host, database)
	default:
		logger.Error.Printf("database engine is not configured")
	}
	return ""
}

func GetMgConnection() *mongo.Database {
	return mg
}
