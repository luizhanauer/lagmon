package database

import (
	"database/sql"
	"fmt"
	"lag-monitor/internal/domain"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteBatcher struct {
	db        *sql.DB
	buffer    []domain.PingResult
	mu        sync.Mutex
	batchSize int
	ticker    *time.Ticker
	quit      chan struct{}
}

func NewSQLiteRepo(dbPath string) (*SQLiteBatcher, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Cria tabela se não existir
	query := `
	CREATE TABLE IF NOT EXISTS pings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		host_id TEXT,
		latency INTEGER,
		jitter INTEGER,
		timestamp DATETIME
	);`
	if _, err := db.Exec(query); err != nil {
		return nil, err
	}

	querySettings := `
	CREATE TABLE IF NOT EXISTS settings (
		key TEXT PRIMARY KEY,
		value TEXT
	);`
	if _, err := db.Exec(querySettings); err != nil {
		return nil, err
	}

	repo := &SQLiteBatcher{
		db:        db,
		batchSize: 100,                             // Grava a cada 100 registros
		ticker:    time.NewTicker(5 * time.Second), // OU a cada 5 segundos
		quit:      make(chan struct{}),
		buffer:    make([]domain.PingResult, 0, 100),
	}

	go repo.flushLoop()
	return repo, nil
}

func (r *SQLiteBatcher) SaveBatch(results []domain.PingResult) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.buffer = append(r.buffer, results...)
	if len(r.buffer) >= r.batchSize {
		go r.flush(r.popBuffer())
	}
	return nil
}

func (r *SQLiteBatcher) flushLoop() {
	for {
		select {
		case <-r.ticker.C:
			r.mu.Lock()
			if len(r.buffer) > 0 {
				data := r.popBuffer()
				go r.flush(data)
			}
			r.mu.Unlock()
		case <-r.quit:
			return
		}
	}
}

func (r *SQLiteBatcher) popBuffer() []domain.PingResult {
	tmp := make([]domain.PingResult, len(r.buffer))
	copy(tmp, r.buffer)
	r.buffer = r.buffer[:0]
	return tmp
}

func (r *SQLiteBatcher) flush(data []domain.PingResult) {
	if len(data) == 0 {
		return
	}

	tx, err := r.db.Begin()
	if err != nil {
		return
	}

	stmt, err := tx.Prepare("INSERT INTO pings(host_id, latency, jitter, timestamp) VALUES(?, ?, ?, ?)")
	if err != nil {
		return
	}
	defer stmt.Close()

	for _, d := range data {
		stmt.Exec(d.HostID, d.Latency, d.Jitter, d.Timestamp)
	}
	tx.Commit()
}

func (r *SQLiteBatcher) Close() error {
	close(r.quit)
	r.ticker.Stop()
	// Flush final
	r.mu.Lock()
	r.flush(r.buffer)
	r.mu.Unlock()
	return r.db.Close()
}

// GetHistory busca registros filtrados por host e data
func (r *SQLiteBatcher) GetHistory(hostID string, start, end time.Time) ([]domain.PingResult, error) {
	query := `
		SELECT host_id, latency, jitter, timestamp 
		FROM pings 
		WHERE host_id = ? AND timestamp BETWEEN ? AND ?
		ORDER BY timestamp ASC`

	rows, err := r.db.Query(query, hostID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.PingResult
	for rows.Next() {
		var res domain.PingResult
		if err := rows.Scan(&res.HostID, &res.Latency, &res.Jitter, &res.Timestamp); err != nil {
			continue
		}
		results = append(results, res)
	}
	return results, nil
}

// CleanOldData remove registros mais antigos que o número de dias especificado
func (r *SQLiteBatcher) CleanOldData(days int) (int64, error) {
	result, err := r.db.Exec("DELETE FROM pings WHERE timestamp < datetime('now', ?)", fmt.Sprintf("-%d days", days))
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Métodos para gerenciar configurações
func (r *SQLiteBatcher) SetSetting(key, value string) error {
	_, err := r.db.Exec("INSERT OR REPLACE INTO settings (key, value) VALUES (?, ?)", key, value)
	return err
}

func (r *SQLiteBatcher) GetSetting(key string) (string, error) {
	var value string
	err := r.db.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	return value, err
}
