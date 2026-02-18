package main

import (
	"context"
	"fmt"
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
}

func NewApp(svc *usecase.MonitorService) *App {
	return &App{service: svc}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Hosts Padrão
	a.service.AddHost(domain.Host{
		ID: "gateway", Name: "Gateway", IP: "192.168.1.1", IsGW: true, Active: true,
	})
	a.service.AddHost(domain.Host{
		ID: "google", Name: "Google DNS", IP: "8.8.8.8", IsGW: false, Active: true,
	})
}

// Métodos JS

func (a *App) AddTarget(ip string, name string) domain.Host {
	host := domain.Host{
		ID: ip, Name: name, IP: ip, IsGW: false, Active: true,
	}
	a.service.AddHost(host)
	return host
}

func (a *App) RemoveTarget(hostID string) {
	a.service.RemoveHost(hostID)
}

func (a *App) GetTargets() []domain.Host {
	return a.service.GetAllHosts()
}

// SetTargetActive controla o ping no backend
func (a *App) SetTargetActive(hostID string, active bool) {
	a.service.ToggleHostStatus(hostID, active)
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

	// RETORNE APENAS O CAMINHO DO ARQUIVO
	return summaryPath, nil
}

func (a *App) OpenPath(path string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		// No Ubuntu, xdg-open decide se abre a pasta ou o arquivo
		cmd = exec.Command("xdg-open", path)
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin":
		cmd = exec.Command("open", path)
	}

	if cmd != nil {
		cmd.Run()
	}
}
