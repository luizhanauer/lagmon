package usecase

import (
	"context"
	"lag-monitor/internal/domain"
	"math"
	"sync"
	"time"
)

// EventEmitter define a função de envio para o frontend
type EventEmitter func(eventName string, data interface{})

// monitorJob representa a tarefa em execução
type monitorJob struct {
	host    domain.Host
	cancel  context.CancelFunc
	lastLat int64
	active  bool // Estado local de execução
}

// MonitorService gerencia os jobs
type MonitorService struct {
	repo   domain.Repository
	pinger domain.Pinger
	emit   EventEmitter

	mu      sync.RWMutex
	targets map[string]*monitorJob
}

// NewMonitorService construtor
func NewMonitorService(r domain.Repository, p domain.Pinger, e EventEmitter) *MonitorService {
	return &MonitorService{
		repo:    r,
		pinger:  p,
		emit:    e,
		targets: make(map[string]*monitorJob),
	}
}

// AddHost adiciona um novo alvo
func (s *MonitorService) AddHost(h domain.Host) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.targets[h.ID]; exists {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Define como ativo por padrão na criação se não especificado
	activeState := true
	if !h.Active && h.ID == "" {
		// Se for um novo host sem config prévia, assume true
		activeState = true
	} else {
		activeState = h.Active
	}

	// Atualiza struct do dominio para refletir o estado
	h.Active = activeState

	job := &monitorJob{
		host:    h,
		cancel:  cancel,
		lastLat: 0,
		active:  activeState,
	}
	s.targets[h.ID] = job

	go s.runLoop(ctx, job)
}

// RemoveHost para o monitoramento e deleta
func (s *MonitorService) RemoveHost(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if job, exists := s.targets[id]; exists {
		job.cancel()
		delete(s.targets, id)
	}
}

// ToggleHostStatus pausa ou retoma o ping
func (s *MonitorService) ToggleHostStatus(id string, active bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if job, exists := s.targets[id]; exists {
		job.active = active
		job.host.Active = active // Atualiza a definição do host também
	}
}

// GetAllHosts retorna lista para UI
func (s *MonitorService) GetAllHosts() []domain.Host {
	s.mu.RLock()
	defer s.mu.RUnlock()

	hosts := make([]domain.Host, 0, len(s.targets))
	for _, job := range s.targets {
		// Garante que o objeto host retornado tenha o estado atualizado
		h := job.host
		h.Active = job.active
		hosts = append(hosts, h)
	}
	return hosts
}

// runLoop executa o ping periodicamente
func (s *MonitorService) runLoop(ctx context.Context, job *monitorJob) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Se estiver pausado, não faz nada mas mantém o loop vivo
			if !job.active {
				continue
			}

			lat, err := s.pinger.Ping(job.host.IP, 900*time.Millisecond)

			res := domain.PingResult{
				HostID:    job.host.ID,
				IP:        job.host.IP,
				Timestamp: time.Now(),
				Loss:      err != nil,
			}

			if err == nil {
				res.Latency = lat
				if job.lastLat > 0 {
					res.Jitter = int64(math.Abs(float64(lat - job.lastLat)))
				}
				job.lastLat = lat
			} else {
				job.lastLat = 0
				res.Latency = 0
				res.Jitter = 0
			}

			s.emit("ping:data", res)
			s.repo.SaveBatch([]domain.PingResult{res})
		}
	}
}
