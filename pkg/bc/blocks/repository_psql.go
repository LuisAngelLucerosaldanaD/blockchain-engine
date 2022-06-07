package blocks

/*// psql estructura de conexión a la BD de postgresql
type psql struct {
	DB   *sqlx.DB
	user *models.User
	TxID string
}

func newBlockPsqlRepository(db *sqlx.DB, user *models.User, txID string) *psql {
	return &psql{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) create(m *Block) error {
	const psqlInsert = `INSERT INTO bc.blocks (id, data, nonce, difficulty, mined_by, mined_at, timestamp, hash, prev_hash) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, created_at, updated_at`
	stmt, err := s.DB.Prepare(psqlInsert)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		m.ID,
		m.Data,
		m.Nonce,
		m.Difficulty,
		m.MinedBy,
		m.MinedAt,
		m.Timestamp,
		m.Hash,
		m.PrevHash,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *psql) update(m *Block) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE bc.blocks SET data = :data, nonce = :nonce, difficulty = :difficulty, mined_by = :mined_by, mined_at = :mined_at, timestamp = :timestamp, hash = :hash, prev_hash = :prev_hash, updated_at = :updated_at WHERE id = :id `
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
	const psqlDelete = `DELETE FROM bc.blocks WHERE id = :id `
	m := Block{ID: id}
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
func (s *psql) getByID(id int64) (*Block, error) {
	const psqlGetByID = `SELECT id , data, nonce, difficulty, mined_by, mined_at, timestamp, hash, prev_hash, created_at, updated_at FROM bc.blocks WHERE id = $1 `
	mdl := Block{}
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
func (s *psql) getAll(limit, offSet *int) ([]*Block, error) {
	var ms []*Block
	const psqlGetAll = `SELECT id , data, nonce, difficulty, mined_by, mined_at, timestamp, hash, prev_hash, created_at, updated_at FROM bc.blocks order by id desc Limit %d offSet %d`

	query := fmt.Sprintf(psqlGetAll, *limit, *offSet)

	err := s.DB.Select(&ms, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

// getBlocksById consulta un bloque de la BD por id
func (s *psql) getBlocksById(id int64) ([]*Block, error) {
	var ms []*Block
	const psqlGetAll = `SELECT id , data, nonce, difficulty, mined_by, mined_at, timestamp, hash, prev_hash, created_at, updated_at FROM bc.blocks where id = %d;`

	query := fmt.Sprintf(psqlGetAll, id)

	err := s.DB.Select(&ms, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

func (s *psql) getHashPrevBlock() (string, error) {
	const queryGetHashPrevBlock = `SELECT hash as prev_hash FROM bc.blocks order by id desc limit 1;`
	var mdl string
	err := s.DB.Get(&mdl, queryGetHashPrevBlock)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return mdl, err
	}
	return mdl, nil
}

func (s *psql) existsBlock() bool {
	var ms []*Block
	const psqlGetAll = `SELECT id , data, nonce, difficulty, mined_by, mined_at, timestamp, hash, prev_hash, created_at, updated_at FROM bc.blocks order by id desc Limit 1`

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		return true
	}
	if ms == nil {
		return false
	}
	return true
}
*/
