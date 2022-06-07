package messages

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
	modelCollection = "messages"
)

type mongodb struct {
	DBMg *mongo.Database
	user *models.User
	TxID string
}

func newMessagesMongodb(dbMg *mongo.Database, user *models.User, txID string) *mongodb {
	return &mongodb{
		DBMg: dbMg,
		user: user,
		TxID: txID,
	}
}

func (mg mongodb) create(m *Messages) error {
	connDB := mg.DBMg

	collection := connDB.Collection(modelCollection)
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	_, err := collection.InsertOne(context.TODO(), &m)
	if err != nil {
		logger.Error.Printf("creando mensages: %v", err)
		return err
	}
	return nil
}

func (mg mongodb) update(m *Messages) error {
	connDB := mg.DBMg
	collection := connDB.Collection(modelCollection)

	MessagesRes, err := mg.getByID(m.ID)
	m.CreatedAt = MessagesRes.CreatedAt
	m.UpdatedAt = time.Now()
	filter := bson.D{{"_id", m.ID}}
	update := bson.M{"$set": m}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		logger.Error.Printf("ejecutando Update mensajes: %v", err)
		return err
	}

	return nil
}

func (mg mongodb) delete(id int) error {
	connDB := mg.DBMg
	collection := connDB.Collection(modelCollection)

	filter := bson.D{{"_id", id}}

	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		logger.Error.Printf("ejecutando Delete mensajes: %v", err)
		return err
	}
	return nil
}

func (mg mongodb) getByID(ID int) (*Messages, error) {
	connDB := mg.DBMg
	collection := connDB.Collection(modelCollection)

	filter := bson.D{{"_id", ID}}
	m := Messages{}

	err := collection.FindOne(context.TODO(), filter).Decode(&m)

	if err != mongo.ErrNoDocuments {
		if err != nil {
			logger.Error.Printf("consultando Get mensajes ByID config: %v", err)
			return nil, err
		}
	}

	return &m, nil
}

func (mg mongodb) getAll() ([]*Messages, error) {
	connDB := mg.DBMg
	collection := connDB.Collection(modelCollection)

	findOptions := options.Find()

	rs, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != mongo.ErrNoDocuments {
		if err != nil {
			logger.Error.Printf("ejecutando GetAllMessages config: %v", err)
			return nil, err
		}
	}

	ms, err := mg.scanRow(rs)
	if err != nil {
		logger.Error.Printf("ejecutando GetAllMessages config: %v", err)
		return nil, err
	}

	return ms, nil
}

func (mg mongodb) scanRow(rs *mongo.Cursor) ([]*Messages, error) {
	var results []*Messages
	for rs.Next(context.TODO()) {
		var elem *Messages
		err := rs.Decode(&elem)
		if err != nil {
			logger.Error.Printf("escaneando el modelo config: %v", err)
			return nil, err
		}
		results = append(results, elem)
	}

	err := rs.Err()
	if err != nil {
		logger.Error.Printf("validando consistencia en el modelo config: %v", err)
		return nil, err
	}

	err = rs.Close(context.TODO())
	if err != nil {
		logger.Error.Printf("couldn't close rs: %v", err)
		return nil, err
	}

	return results, nil
}
