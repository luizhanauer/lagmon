package domain

import "time"

// PingResult representa um único ponto de dados
type PingResult struct {
	HostID    string    `json:"hostId"`
	IP        string    `json:"ip"`
	Latency   int64     `json:"latency"` // em microsegundos
	Jitter    int64     `json:"jitter"`  // em microsegundos
	Timestamp time.Time `json:"timestamp"`
	Loss      bool      `json:"loss"`
}

// Host define um alvo para monitoramento
type Host struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	IP     string `json:"ip"`
	IsGW   bool   `json:"isGateway"`
	Active bool   `json:"active"` // Novo campo para controle de estado
}

// Repository define como salvamos os dados
type Repository interface {
	SaveBatch(results []PingResult) error
	Close() error
}

// Pinger define como executamos o ping
type Pinger interface {
	Ping(ip string, timeout time.Duration) (int64, error)
}
