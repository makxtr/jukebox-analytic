package main

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

type sqliteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(dbPath string) (IRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	repo := &sqliteRepository{db: db}
	if err := repo.migrate(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return repo, nil
}

func (r *sqliteRepository) migrate() error {
	migrationSQL, err := migrationFS.ReadFile("migrations/001_init.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	_, err = r.db.Exec(string(migrationSQL))
	return err
}

func (r *sqliteRepository) GetTrackByID(id int) (*Track, error) {
	row := r.db.QueryRow("SELECT id, title, artist, price FROM tracks WHERE id = ?", id)
	var t Track
	if err := row.Scan(&t.ID, &t.Title, &t.Artist, &t.Price); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("track with id %d not found", id)
		}
		return nil, err
	}
	return &t, nil
}

func (r *sqliteRepository) UpdateTrackPrice(id int, newPrice float64) error {
	res, err := r.db.Exec("UPDATE tracks SET price = ? WHERE id = ?", newPrice, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("track with id %d not found", id)
	}
	return nil
}

func (r *sqliteRepository) CreateLog(log PlaybackLog) error {
	_, err := r.db.Exec("INSERT INTO playback_logs (track_id, played_at, amount_paid) VALUES (?, ?, ?)",
		log.TrackID, log.PlayedAt, log.AmountPaid)
	return err
}

func (r *sqliteRepository) GetAllLogs() []PlaybackLog {
	rows, err := r.db.Query("SELECT id, track_id, played_at, amount_paid FROM playback_logs")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var logs []PlaybackLog
	for rows.Next() {
		var l PlaybackLog
		if err := rows.Scan(&l.ID, &l.TrackID, &l.PlayedAt, &l.AmountPaid); err == nil {
			logs = append(logs, l)
		}
	}
	return logs
}

func (r *sqliteRepository) GetTopTracks(limit int) ([]TopTrackStat, error) {
	query := `
		SELECT t.title, COUNT(l.id) as play_count
		FROM playback_logs l
		JOIN tracks t ON l.track_id = t.id
		GROUP BY t.id
		ORDER BY play_count DESC
		LIMIT ?
	`
	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []TopTrackStat
	for rows.Next() {
		var s TopTrackStat
		if err := rows.Scan(&s.Title, &s.Count); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}
