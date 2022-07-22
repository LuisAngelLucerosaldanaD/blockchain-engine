package blocks_tmp

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
	modelCollection = "blocks_temp"
)

// psql estructura de conexión a la BD de postgresql
type mongodb struct {
	DBMg *mongo.Database
	user *models.User
	TxID string
}

func newBlockTempMongodbRepository(dbMg *mongo.Database, user *models.User, txID string) *mongodb {
	return &mongodb{
		DBMg: dbMg,
		user: user,
		TxID: txID,
	}
}

func (mg mongodb) create(m *BlockTmp) error {
	connDB := mg.DBMg

	collection := connDB.Collection(modelCollection)
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	nBlocks, err := collection.CountDocuments(context.TODO(), bson.D{{}})
	if err != nil {
		logger.Error.Printf("Error creando el block, no id value, error: %v", err)
		return err
	}

	m.ID = nBlocks + 1
	_, err = collection.InsertOne(context.TODO(), &m)
	if err != nil {
		logger.Error.Printf("Error creando el bloque temporal, error: %v", err)
		return err
	}
	return nil
}

func (mg mongodb) update(m *BlockTmp) error {
	connDB := mg.DBMg
	collection := connDB.Collection(modelCollection)

	block, err := mg.getByID(m.ID)
	m.Timestamp = block.Timestamp
	m.CreatedAt = block.CreatedAt
	m.UpdatedAt = time.Now()
	filter := bson.D{{"_id", m.ID}}
	update := bson.M{"$set": m}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		logger.Error.Printf("Error actualizando el bloque temporal, error: %v", err)
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
		logger.Error.Printf("error eliminando el bloque temporal, error: %v", err)
		return err
	}
	return nil
}

func (mg mongodb) getByID(ID int64) (*BlockTmp, error) {
	connDB := mg.DBMg
	collection := connDB.Collection(modelCollection)

	filter := bson.D{{"_id", ID}}
	m := BlockTmp{}

	err := collection.FindOne(context.TODO(), filter).Decode(&m)
	if err != mongo.ErrNoDocuments {
		if err != nil {
			logger.Error.Printf("Error trayendo el bloque temporal por ID, error: %v", err)
			return nil, err
		}
	}

	return &m, nil
}

func (mg mongodb) getAll() ([]*BlockTmp, error) {
	connDB := mg.DBMg
	collection := connDB.Collection(modelCollection)
	findOptions := options.Find()

	rs, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != mongo.ErrNoDocuments {
		if err != nil {
			logger.Error.Printf("Error trayendo todos los bloques temporales, error: %v", err)
			return nil, err
		}
	}

	ms, err := mg.scanRow(rs)
	if err != nil {
		logger.Error.Printf("Error trayendo todos los bloques temporales, error: %v", err)
		return nil, err
	}

	return ms, nil
}

func (mg mongodb) scanRow(rs *mongo.Cursor) ([]*BlockTmp, error) {
	var results []*BlockTmp
	for rs.Next(context.TODO()) {
		var elem *BlockTmp
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

// getBlockUnCommit
func (mg mongodb) getBlockUnCommit() (*BlockTmp, error) {
	mdl := BlockTmp{}

	connDB := mg.DBMg
	collection := connDB.Collection(modelCollection)

	findOptions := options.FindOneOptions{}
	findOptions.SetSort(bson.D{{"id", -1}})

	err := collection.FindOne(context.TODO(), bson.D{{"status", 1}}, &findOptions).Decode(&mdl)
	if err != mongo.ErrNoDocuments {
		if err != nil {
			logger.Error.Printf("Error trayendo el bloque con estado 1, error: %v", err)
			return nil, err
		}
	}

	return &mdl, nil
}

// getBlockUnCommit
func (mg mongodb) getBlockTwoCommit() (*BlockTmp, error) {
	mdl := BlockTmp{}

	connDB := mg.DBMg
	collection := connDB.Collection(modelCollection)

	findOptions := options.FindOneOptions{}
	findOptions.SetSort(bson.D{{"id", -1}})

	err := collection.FindOne(context.TODO(), bson.D{{"status", 2}}, &findOptions).Decode(&mdl)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		logger.Error.Printf("Error trayendo el bloque con estado 2, error: %v", err)
		return nil, err
	}
	return &mdl, nil
}

func (mg mongodb) GetCountTransactionByID(block int64) int64 {
	connDB := mg.DBMg
	collection := connDB.Collection("transactions")

	totalTransaction, err := collection.CountDocuments(context.TODO(), bson.D{{"block", block}})
	if err != mongo.ErrNoDocuments {
		if err != nil {
			logger.Error.Printf("Error trayendo todos los bloques temporales, error: %v", err)
			return 0
		}
	}

	return totalTransaction
}
