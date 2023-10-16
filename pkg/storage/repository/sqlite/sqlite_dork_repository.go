package repository_sqlite

import (
	"database/sql"
	"magma/pkg/dto"
	"magma/pkg/storage/repository"
)

type sqliteDorkRepository struct {
	DB *sql.DB
}

// UpdateDork implements repository.DorkRepository.
func (s *sqliteDorkRepository) UpdateDork(dork dto.DorkDTO) error {
	query := "UPDATE dorks SET score = ? WHERE id = ?"
	statement, err := s.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(dork.Score, dork.ID)
	if err != nil {
		return err
	}

	return nil
}

// SaveDork implements repository.DorkRepository.
func (s *sqliteDorkRepository) SaveDork(dork dto.DorkDTO) error {
	query := "INSERT INTO dorks (dork, score) VALUES (?, ?)"
	statement, err := s.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(dork.Dork, dork.Score)
	if err != nil {
		return err
	}

	return nil
}

// GetDorks implements repository.DorkRepository.
func (s *sqliteDorkRepository) GetDorks() ([]dto.DorkDTO, error) {
	query := "SELECT id, dork, score FROM dorks ORDER BY score DESC"

	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}

	var dorks []dto.DorkDTO
	for rows.Next() {
		var dork dto.DorkDTO
		if err := rows.Scan(&dork.ID, &dork.Dork, &dork.Score); err != nil {
			return nil, err
		}
		dorks = append(dorks, dork)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dorks, nil
}

func NewSqliteDorkRepository(db *sql.DB) repository.DorkRepository {
	return &sqliteDorkRepository{
		DB: db,
	}
}
