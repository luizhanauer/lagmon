package database

import (
	"database/sql"
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
