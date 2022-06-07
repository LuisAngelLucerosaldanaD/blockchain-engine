package blocks_tmp

import (
	"fmt"
	"time"

	"bjungle/blockchain-engine/internal/logger"
	"bjungle/blockchain-engine/internal/models"
)

type PortsServerBlockTmp interface {
	CreateBlockTmp(status int, timestamp time.Time) (*BlockTmp, int, error)
	UpdateBlockTmp(id int64, status int) (*BlockTmp, int, error)
	DeleteBlockTmp(id int64) (int, error)
	GetBlockTmpByID(id int64) (*BlockTmp, int, error)
	GetAllBlockTmp() ([]*BlockTmp, error)
	GetBlockUnCommit() (*BlockTmp, error)
	GetBlockTwoCommit() (*BlockTmp, error)
	/*MustCloseBlock(lifeBlock time.Time, block int64) bool*/
}

type service struct {
	repository ServicesBlockTmpRepository
	user       *models.User
	txID       string
}

func NewBlockTmpService(repository ServicesBlockTmpRepository, user *models.User, TxID string) PortsServerBlockTmp {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateBlockTmp(status int, timestamp time.Time) (*BlockTmp, int, error) {
	m := NewCreateBlockTmp(status, timestamp)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create BlockTmp :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateBlockTmp(id int64, status int) (*BlockTmp, int, error) {
	m := NewBlockTmp(id, status)
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return m, 15, fmt.Errorf("id is required")
	}
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update BlockTmp :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteBlockTmp(id int64) (int, error) {
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

func (s *service) GetBlockTmpByID(id int64) (*BlockTmp, int, error) {
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

func (s *service) GetAllBlockTmp() ([]*BlockTmp, error) {
	return s.repository.getAll()
}

func (s *service) GetBlockUnCommit() (*BlockTmp, error) {
	return s.repository.getBlockUnCommit()
}

func (s *service) GetBlockTwoCommit() (*BlockTmp, error) {
	return s.repository.getBlockTwoCommit()
}

/*func (s *service) MustCloseBlock(TtlBlock time.Time, block int64) bool {
	c := env.NewConfiguration()
	transactions := s.repository.GetCountTransactionByID(block)

	if transactions == 0 {
		return false
	}

	if transactions >= int64(c.App.MaxTransactionsBlock) {
		return true
	}

	lifeBlock := time.Now().Sub(TtlBlock).Seconds()

	if int(lifeBlock) > (c.App.TtlBlock * 1000) {
		return true
	}
	return false
}*/
