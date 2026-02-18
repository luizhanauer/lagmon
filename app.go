package main

import (
	"context"
	"fmt"
	"lag-monitor/internal/config"
	"lag-monitor/internal/domain"
	"lag-monitor/internal/usecase"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

type App struct {
	ctx     context.Context
	service *usecase.MonitorService
	repo    domain.Repository
	cfg     *config.ConfigManager
}

func NewApp(svc *usecase.MonitorService, r domain.Repository, cfg *config.ConfigManager) *App {
	return &App{
		service: svc,
		repo:    r,
		cfg:     cfg,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	a.service.StartRetentionPolicy()

	if err := a.cfg.Load(); err != nil {
		fmt.Println("Erro ao carregar config:", err)
	}

	for _, target := range a.cfg.Data.Targets {
		a.service.AddHost(target)
	}
}

// --- NOVOS MÉTODOS PARA PERSISTÊNCIA NO SETTINGS.JSON ---

// GetConfig exporta a configuração atual do arquivo para o frontend
func (a *App) GetConfig() config.AppConfig {
	return a.cfg.Data
}

// UpdateConfig recebe a configuração do frontend e salva no settings.json
func (a *App) UpdateConfig(newCfg config.AppConfig) error {
	a.cfg.UpdateConfig(newCfg)
	return a.cfg.Save()
}

// --- MÉTODOS JS EXISTENTES ---

func (a *App) AddTarget(ip string, name string) domain.Host {
	host := domain.Host{
		ID: ip, Name: name, IP: ip, IsGW: false, Active: true,
	}

	a.service.AddHost(host)
	a.cfg.AddTarget(host)
	return host
}

func (a *App) RemoveTarget(hostID string) {
	a.service.RemoveHost(hostID)
	a.cfg.RemoveTarget(hostID)
}

func (a *App) GetTargets() []domain.Host {
	return a.service.GetAllHosts()
}

func (a *App) SetTargetActive(hostID string, active bool) {
	a.service.ToggleHostStatus(hostID, active)
	a.cfg.UpdateTargetStatus(hostID, active)
}

func (a *App) GetReport(hostID string, startStr, endStr string) (string, error) {
	layout := "2006-01-02T15:04"
	start, _ := time.Parse(layout, startStr)
	end, _ := time.Parse(layout, endStr)

	summary, raw, err := a.service.GenerateDualReport(hostID, start, end)
	if err != nil {
		return "", err
	}

	home, _ := os.UserHomeDir()
	timestamp := time.Now().Unix()

	summaryPath := filepath.Join(home, "Downloads", fmt.Sprintf("RESUMO-%s-%d.txt", hostID, timestamp))
	os.WriteFile(summaryPath, []byte(summary), 0644)

	rawPath := filepath.Join(home, "Downloads", fmt.Sprintf("DADOS-TECNICOS-%s-%d.csv", hostID, timestamp))
	os.WriteFile(rawPath, []byte(raw), 0644)

	return summaryPath, nil
}

func (a *App) OpenPath(path string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", path)
	case "windows":
		cmd = exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", path)
	case "darwin":
		cmd = exec.Command("open", path)
	}

	if cmd != nil {
		cmd.Run()
	}
}

func (a *App) GetDiagramConfig() config.NetworkDiagramConfig {
	return a.cfg.Data.NetworkDiagram
}

func (a *App) SetTargetDiagramVisibility(hostID string, show bool) {
	// 1. Atualiza o estado no arquivo settings.json
	a.cfg.UpdateTargetDiagramVisibility(hostID, show)

	// 2. Se você precisar que o Service saiba disso em tempo real:
	// a.service.ToggleDiagramStatus(hostID, show)
}
