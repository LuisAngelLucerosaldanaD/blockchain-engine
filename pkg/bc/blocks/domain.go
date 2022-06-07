package blocks

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Block estructura de Block
type Block struct {
	ID                 int64     `json:"id" db:"id" bson:"_id" valid:"-"`
	Data               string    `json:"data" db:"data" bson:"data" valid:"required"`
	Nonce              int64     `json:"nonce" db:"nonce" bson:"nonce" valid:"required"`
	Difficulty         int       `json:"difficulty" db:"difficulty" bson:"difficulty" valid:"required"`
	MinedBy            string    `json:"mined_by" db:"mined_by" bson:"mined_by" valid:"required"`
	MinedAt            time.Time `json:"mined_at" db:"mined_at" bson:"mined_at" valid:"required"`
	Timestamp          time.Time `json:"timestamp" db:"timestamp" bson:"timestamp" valid:"required"`
	Hash               string    `json:"hash" db:"hash" bson:"hash" valid:"required"`
	PrevHash           string    `json:"prev_hash" db:"prev_hash" bson:"prev_hash" valid:"required"`
	StatusId           int       `json:"status_id" db:"status_id" bson:"status_id" valid:"required"`
	IdUser             string    `json:"id_user" db:"id_user" bson:"id_user" valid:"required"`
	LastValidationDate string    `json:"last_validation_date" db:"last_validation_date" bson:"last_validation_date" valid:"required"`
	CreatedAt          time.Time `json:"created_at"  db:"created_at" bson:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at" bson:"updated_at"`
}

func NewBlock(id int64, data string, nonce int64, difficulty int, minedBy string, minedAt time.Time, timestamp time.Time, hash string, prevHash string) *Block {
	return &Block{
		ID:                 id,
		Data:               data,
		Nonce:              nonce,
		Difficulty:         difficulty,
		MinedBy:            minedBy,
		MinedAt:            minedAt,
		Timestamp:          timestamp,
		Hash:               hash,
		PrevHash:           prevHash,
		StatusId:           12,
		IdUser:             minedBy,
		LastValidationDate: time.Now().String(),
	}
}

func (m *Block) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
