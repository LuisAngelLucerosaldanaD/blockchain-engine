package blocks

import (
	"bjungle/blockchain-engine/internal/logger"
	"bjungle/blockchain-engine/internal/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	modelCollection = "blocks"
)

// psql estructura de conexión a la BD de postgresql
type mongodb struct {
	DBMg *mongo.Database
	user *models.User
	TxID string
}

func newBlockMongodbRepository(dbMg *mongo.Database, user *models.User, txID string) *mongodb {
	return &mongodb{
		DBMg: dbMg,
		user: user,
		TxID: txID,
	}
}

func (mg mongodb) create(m *Block) error {
	connDB := mg.DBMg

	collection := connDB.Collection(modelCollection)
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	_, err := collection.InsertOne(context.TODO(), &m)
	if err != nil {
		logger.Error.Printf("Error creando el bloque, error: %v", err)
		return err
	}
	return nil
}

func (mg mongodb) update(m *Block) error {
	connDB := mg.DBMg
	collection := connDB.Collection(modelCollection)

	block, err := mg.getByID(m.ID)
	m.CreatedAt = block.CreatedAt
	m.UpdatedAt = time.Now()
	filter := bson.D{{"_id", m.ID}}
	update := bson.M{"$set": m}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		logger.Error.Printf("Error actualizando el bloque, error: %v", err)
		return err
	}

	return nil
}

func (mg mongodb) delete(id int64) error {
	connDB := mg.DBMg
	collection := connDB.Collection(modelCollection)

	filter := bson.D{{"_id", id}}

	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		logger.Error.Printf("ejecutando Delete: %v", err)
		return err
	}
	return nil
}

func (mg mongodb) getByID(ID int64) (*Block, error) {
	connDB := mg.DBMg
	collection := connDB.Collection(modelCollection)

	filter := bson.D{{"_id", ID}}
	m := Block{}

	err := collection.FindOne(context.TODO(), filter).Decode(&m)

	if err != mongo.ErrNoDocuments {
		if err != nil {
			logger.Error.Printf("Error trayendo el bloque por ID, error: %v", err)
			return nil, err
		}
	}

	return &m, nil
}

func (mg mongodb) getAll(limit, offSet *int64) ([]*Block, error) {
	connDB := mg.DBMg
	collection := connDB.Collection(modelCollection)

	findOptions := options.Find()
	findOptions.SetLimit(*limit)
	findOptions.SetSkip(*offSet)

	rs, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != mongo.ErrNoDocuments {
		if err != nil {
			logger.Error.Printf("Error trayendo todos los bloques, error: %v", err)
			return nil, err
		}
	}

	ms, err := mg.scanRow(rs)
	if err != nil {
		logger.Error.Printf("Error trayendo todos los bloques, error: %v", err)
		return nil, err
	}

	return ms, nil
}

func (mg mongodb) getBlocksById(id int64) (*Block, error) {
	connDB := mg.DBMg
	collection := connDB.Collection(modelCollection)

	filter := bson.D{{"_id", id}}
	m := Block{}

	err := collection.FindOne(context.TODO(), filter).Decode(&m)
	if err != mongo.ErrNoDocuments {
		if err != nil {
			logger.Error.Printf("Error trayendo el bloque por ID, error: %v", err)
			return nil, err
		}
	}

	return &m, nil
}

func (mg mongodb) getHashPrevBlock() (string, error) {
	var ms *Block

	connDB := mg.DBMg
	collection := connDB.Collection(modelCollection)

	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{"_id", -1}})

	err := collection.FindOne(context.TODO(), bson.D{}, findOptions).Decode(&ms)
	if err != mongo.ErrNoDocuments {
		if err != nil {
			logger.Error.Printf("Error trayendo todos los bloques, error: %v", err)
			return "", err
		}
	}

	return ms.Hash, nil
}

func (mg mongodb) existsBlock() bool {
	var ms *Block

	connDB := mg.DBMg
	collection := connDB.Collection(modelCollection)

	err := collection.FindOne(context.TODO(), bson.D{{"id", -1}}).Decode(&ms)
	if err != mongo.ErrNoDocuments {
		if err != nil {
			logger.Error.Printf("Error trayendo todos los bloques, error: %v", err)
			return false
		}
	}

	if ms == nil {
		return false
	}
	return true
}

func (mg mongodb) scanRow(rs *mongo.Cursor) ([]*Block, error) {
	var results []*Block
	for rs.Next(context.TODO()) {
		var elem *Block
		err := rs.Decode(&elem)
		if err != nil {
			logger.Error.Printf("escaneando el modelo block: %v", err)
			return nil, err
		}
		results = append(results, elem)
	}

	err := rs.Err()
	if err != nil {
		logger.Error.Printf("validando consistencia en el modelo block: %v", err)
		return nil, err
	}

	err = rs.Close(context.TODO())
	if err != nil {
		logger.Error.Printf("couldn't close rs: %v", err)
		return nil, err
	}

	return results, nil
}
