package main

import (
	"context"
	"lag-monitor/internal/domain"
	"lag-monitor/internal/usecase"
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
