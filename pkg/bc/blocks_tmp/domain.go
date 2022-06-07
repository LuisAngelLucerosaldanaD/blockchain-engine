package blocks_tmp

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Model estructura de BlockTmp
type BlockTmp struct {
	ID        int64     `json:"id" db:"id" bson:"_id" valid:"-"`
	Status    int       `json:"status" db:"status" bson:"status"`
	Timestamp time.Time `json:"timestamp" db:"timestamp" bson:"timestamp"`
	CreatedAt time.Time `json:"created_at" db:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" bson:"updated_at"`
}

func NewBlockTmp(id int64, status int) *BlockTmp {
	return &BlockTmp{
		ID:     id,
		Status: status,
	}
}

func NewCreateBlockTmp(status int, timestamp time.Time) *BlockTmp {
	return &BlockTmp{
		Status:    status,
		Timestamp: timestamp,
	}
}

func (m *BlockTmp) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
