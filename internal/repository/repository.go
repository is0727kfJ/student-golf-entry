package repository

import (
	"database/sql"
	"errors"

	"github.com/is0727kfJ/student-golf-entry/internal/models"
)

// ITournamentRepository は倉庫番が必ず守るべきルール（インターフェース）
type ITournamentRepository interface {
	GetAll() ([]models.Tournament, error)
	CreateEntryTx(tournamentID, userID string) error
}

type tournamentRepository struct {
	db *sql.DB
}

func NewTournamentRepository(db *sql.DB) ITournamentRepository {
	return &tournamentRepository{db: db}
}

func (r *tournamentRepository) GetAll() ([]models.Tournament, error) {
	query := `SELECT id, title, capacity, current_entries, status FROM tournaments ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tournaments []models.Tournament
	for rows.Next() {
		var t models.Tournament
		if err := rows.Scan(&t.ID, &t.Title, &t.Capacity, &t.CurrentEntries, &t.Status); err != nil {
			return nil, err
		}
		tournaments = append(tournaments, t)
	}
	return tournaments, nil
}

// トランザクション処理（排他制御）はデータベース特有の処理なのでRepositoryに任せる
func (r *tournamentRepository) CreateEntryTx(tournamentID, userID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	updateQuery := `UPDATE tournaments SET current_entries = current_entries + 1 WHERE id = $1 AND current_entries < capacity`
	res, err := tx.Exec(updateQuery, tournamentID)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("capacity_full_or_not_found")
	}

	insertQuery := `INSERT INTO entries (tournament_id, user_id, status) VALUES ($1, $2, 'RESERVED')`
	if _, err = tx.Exec(insertQuery, tournamentID, userID); err != nil {
		return errors.New("duplicate_entry")
	}

	return tx.Commit()
}
