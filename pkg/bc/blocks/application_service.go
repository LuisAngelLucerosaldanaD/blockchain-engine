package blocks

import (
	"fmt"
	"time"

	"bjungle/blockchain-engine/internal/logger"
	"bjungle/blockchain-engine/internal/models"
)

type PortsServerBlock interface {
	CreateBlock(id int64, data string, nonce int64, difficulty int, minedBy string, minedAt time.Time, timestamp time.Time, hash string, prevHash string) (*Block, int, error)
	UpdateBlock(id int64, data string, nonce int64, difficulty int, minedBy string, minedAt time.Time, timestamp time.Time, hash string, prevHash string) (*Block, int, error)
	DeleteBlock(id int64) (int, error)
	GetBlockByID(id int64) (*Block, int, error)
	GetAllBlock(limit, offSet *int64) ([]*Block, error)
	GetHashPrevBlock() (string, error)
	GetBlocksById(id int64) (*Block, error)
	ExistsBlocks() bool
}

type service struct {
	repository ServicesBlockRepository
	user       *models.User
	txID       string
}

func NewBlockService(repository ServicesBlockRepository, user *models.User, TxID string) PortsServerBlock {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateBlock(id int64, data string, nonce int64, difficulty int, minedBy string, minedAt time.Time, timestamp time.Time, hash string, prevHash string) (*Block, int, error) {
	m := NewBlock(id, data, nonce, difficulty, minedBy, minedAt, timestamp, hash, prevHash)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create Block :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateBlock(id int64, data string, nonce int64, difficulty int, minedBy string, minedAt time.Time, timestamp time.Time, hash string, prevHash string) (*Block, int, error) {
	m := NewBlock(id, data, nonce, difficulty, minedBy, minedAt, timestamp, hash, prevHash)
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return m, 15, fmt.Errorf("id is required")
	}
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update Block :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteBlock(id int64) (int, error) {
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return 15, fmt.Errorf("id is required")
	}

	if err := s.repository.delete(id); err != nil {
		if err.Error() == "ecatch:108" {
			return 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update row:", err)
		return 20, err
	}
	return 28, nil
}

func (s *service) GetBlockByID(id int64) (*Block, int, error) {
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return nil, 15, fmt.Errorf("id is required")
	}
	m, err := s.repository.getByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetAllBlock(limit, offSet *int64) ([]*Block, error) {
	return s.repository.getAll(limit, offSet)
}

func (s *service) GetBlocksById(id int64) (*Block, error) {
	return s.repository.getBlocksById(id)
}

func (s *service) GetHashPrevBlock() (string, error) {
	return s.repository.getHashPrevBlock()
}

func (s *service) ExistsBlocks() bool {
	return s.repository.existsBlock()
}
