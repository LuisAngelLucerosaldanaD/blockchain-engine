package blocks_tmp

/*// psql estructura de conexión a la BD de postgresql
type psql struct {
	DB   *sqlx.DB
	user *models.User
	TxID string
}

func newBlockTmpPsqlRepository(db *sqlx.DB, user *models.User, txID string) *psql {
	return &psql{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) create(m *BlockTmp) error {
	const psqlInsert = `INSERT INTO bc.blocks_tmp (status, timestamp) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	stmt, err := s.DB.Prepare(psqlInsert)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		m.Status,
		m.Timestamp,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *psql) update(m *BlockTmp) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE bc.blocks_tmp SET status = :status, updated_at = :updated_at WHERE id = :id `
	rs, err := s.DB.NamedExec(psqlUpdate, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Delete elimina un registro de la BD
func (s *psql) delete(id int64) error {
	const psqlDelete = `DELETE FROM bc.blocks_tmp WHERE id = :id `
	m := BlockTmp{ID: id}
	rs, err := s.DB.NamedExec(psqlDelete, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// GetByID consulta un registro por su ID
func (s *psql) getByID(id int64) (*BlockTmp, error) {
	const psqlGetByID = `SELECT id , status, timestamp, created_at, updated_at FROM bc.blocks_tmp WHERE id = $1 `
	mdl := BlockTmp{}
	err := s.DB.Get(&mdl, psqlGetByID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

// GetAll consulta todos los registros de la BD
func (s *psql) getAll() ([]*BlockTmp, error) {
	var ms []*BlockTmp
	const psqlGetAll = ` SELECT id , status, timestamp, created_at, updated_at FROM bc.blocks_tmp `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

// getBlockUnCommit
func (s *psql) getBlockUnCommit() (*BlockTmp, error) {
	const queryGetBlockUnStatus = `select bt.id, bt.status, bt."timestamp" from bc.blocks_tmp bt where bt.status = 1 order by bt.id desc limit 1;`
	mdl := BlockTmp{}
	err := s.DB.Get(&mdl, queryGetBlockUnStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

// getBlockUnCommit
func (s *psql) getBlockTwoCommit() (*BlockTmp, error) {
	const queryGetBlockTwoCommit = `select bt.id, bt.status, bt."timestamp" from bc.blocks_tmp bt where bt.status = 2 order by id asc limit 1;`
	mdl := BlockTmp{}
	err := s.DB.Get(&mdl, queryGetBlockTwoCommit)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

func (s *psql) GetCountTransactionByID(block int64) int {
	const queryGetCountTransactionByID = `select count(*) from bc.transactions t where t.block = $1;`
	var totalTransaction int
	err := s.DB.Get(&totalTransaction, queryGetCountTransactionByID, block)
	if err != nil {
		if err == sql.ErrNoRows {
			return totalTransaction
		}
		return totalTransaction
	}
	return totalTransaction
}
*/
