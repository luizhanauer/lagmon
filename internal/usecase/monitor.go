package usecase

import (
	"context"
	"fmt"
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
			if !job.active {
				continue
			}

			// Executa o Ping
			lat, err := s.pinger.Ping(job.host.IP, 900*time.Millisecond)

			res := domain.PingResult{
				HostID:    job.host.ID,
				IP:        job.host.IP,
				Timestamp: time.Now(),
				Loss:      err != nil, // Se houver erro, Loss é true
			}

			if err == nil {
				res.Latency = lat
				// Cálculo de Jitter: Só calcula se o pacote anterior E o atual forem bem sucedidos
				if job.lastLat > 0 {
					res.Jitter = int64(math.Abs(float64(lat - job.lastLat)))
				}
				job.lastLat = lat
			} else {
				// IMPORTANTE: Em caso de LOSS, resetamos o lastLat para não calcular
				// jitter inválido no próximo ping bem sucedido.
				job.lastLat = 0
				res.Latency = 0 // Latência 0 indica tecnicamente indisponível
				res.Jitter = 0
			}

			// Emite para o frontend e salva no banco
			s.emit("ping:data", res)
			s.repo.SaveBatch([]domain.PingResult{res})
		}
	}
}

func (s *MonitorService) GenerateReport(hostID string, start, end time.Time) (string, error) {
	data, err := s.repo.GetHistory(hostID, start, end)
	if err != nil {
		return "", err
	}

	report := "RELATÓRIO DE LATÊNCIA - LAG MONITOR\n"
	report += "Target: " + hostID + "\n"
	report += "Período: " + start.Format("02/01/2006 15:04") + " até " + end.Format("02/01/2006 15:04") + "\n"
	report += "--------------------------------------------------\n"

	for _, d := range data {
		report += d.Timestamp.Format("15:04:05") + " | Lat: " +
			fmt.Sprintf("%vms", d.Latency/1000) + " | Jitter: " +
			fmt.Sprintf("%vms", d.Jitter/1000) + "\n"
	}

	return report, nil
}

func (s *MonitorService) GenerateDualReport(hostID string, start, end time.Time) (string, string, error) {
	data, err := s.repo.GetHistory(hostID, start, end)
	if err != nil {
		return "", "", err
	}

	if len(data) == 0 {
		return "", "", fmt.Errorf("sem dados no período")
	}

	// --- CÁLCULOS PARA O RESUMO (LEIGO) ---
	var totalLat int64
	var maxLat int64
	var minLat int64 = 999999
	var lossCount int

	for _, d := range data {
		if d.Loss {
			lossCount++
			continue
		}
		latMs := d.Latency / 1000
		totalLat += latMs
		if latMs > maxLat {
			maxLat = latMs
		}
		if latMs < minLat {
			minLat = latMs
		}
	}

	count := int64(len(data) - lossCount)
	avgLat := int64(0)
	if count > 0 {
		avgLat = totalLat / count
	}
	lossPct := (float64(lossCount) / float64(len(data))) * 100

	// --- 1. RELATÓRIO AMIGÁVEL (RESUMO) ---
	summary := fmt.Sprintf("=== RELATÓRIO DE QUALIDADE DE INTERNET ===\n")
	summary += fmt.Sprintf("Destino: %s\n", hostID)
	summary += fmt.Sprintf("Período: %s até %s\n", start.Format("02/01 15:04"), end.Format("02/01 15:04"))
	summary += "------------------------------------------\n"
	summary += fmt.Sprintf("Média de Atraso (Latência): %dms\n", avgLat)

	status := "EXCELENTE"
	if avgLat > 100 || lossPct > 2 {
		status = "INSTÁVEL"
	}
	if avgLat > 200 || lossPct > 5 {
		status = "CRÍTICO / RUIM"
	}

	summary += fmt.Sprintf("Status da Conexão: %s\n", status)
	summary += fmt.Sprintf("Perda de Sinal: %.1f%%\n", lossPct)
	summary += "------------------------------------------\n"
	summary += "DICA: Valores acima de 100ms ou perdas de sinal podem causar travamentos em vídeos e jogos.\n"

	// --- 2. DADOS BRUTOS (TÉCNICO) ---
	raw := "TIMESTAMP;LATENCY_MS;JITTER_MS;LOSS\n"
	for _, d := range data {
		raw += fmt.Sprintf("%s;%d;%d;%v\n",
			d.Timestamp.Format("2006-01-02 15:04:05"),
			d.Latency/1000,
			d.Jitter/1000,
			d.Loss)
	}

	return summary, raw, nil
}

func (s *MonitorService) StartRetentionPolicy() {
	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		// O loop 'range' aguarda o canal enviar o sinal a cada 1 hora
		for range ticker.C {
			s.runCleanup()
		}
	}()
}

// Extraído para uma função menor (Clean Code/Single Responsibility)
func (s *MonitorService) runCleanup() {
	daysStr, err := s.repo.GetSetting("retention_days")
	if err != nil {
		return
	}

	var days int
	if _, err := fmt.Sscanf(daysStr, "%d", &days); err != nil {
		return
	}

	s.repo.CleanOldData(days)
}
