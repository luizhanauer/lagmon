package domain

import "time"

// PingResult representa um único ponto de dados
type PingResult struct {
	HostID    string    `json:"hostId"`
	IP        string    `json:"ip"`
	Latency   int64     `json:"latency"` // em microsegundos
	Jitter    int64     `json:"jitter"`  // em microsegundos
	Loss      bool      `json:"loss"`
	Timestamp time.Time `json:"timestamp"`
}

// Host define um alvo para monitoramento
type Host struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	IP            string `json:"ip"`
	IsGW          bool   `json:"isGateway"`
	Active        bool   `json:"active"`
	ShowInDiagram bool   `json:"showInDiagram"` // Novo campo
}

// Repository define como salvamos os dados
type Repository interface {
	SaveBatch(results []PingResult) error
	Close() error
	GetHistory(hostID string, start, end time.Time) ([]PingResult, error)
	CleanOldData(days int) (int64, error)
	SetSetting(key, value string) error
	GetSetting(key string) (string, error)
}

// Pinger define como executamos o ping
type Pinger interface {
	Ping(ip string, timeout time.Duration) (int64, error)
}
