package main

import (
	"embed"
	"lag-monitor/internal/infra/database"
	"lag-monitor/internal/infra/network"
	"lag-monitor/internal/usecase"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 1. Infra
	repo, err := database.NewSQLiteRepo("./lagmonitor.db")
	if err != nil {
		log.Fatal(err)
	}

	pinger := network.NewPinger()

	// 2. Serviço (com closure para emitir eventos)
	// Precisamos de uma variável temporária para o contexto do Wails,
	// mas como o emit acontece DEPOIS do start, usamos um canal ou wrapper.
	// Solução simples: passar a função que busca o contexto global do App.

	var app *App // Declaração antecipada

	emitter := func(event string, data interface{}) {
		if app != nil && app.ctx != nil {
			runtime.EventsEmit(app.ctx, event, data)
		}
	}

	service := usecase.NewMonitorService(repo, pinger, emitter)
	app = NewApp(service, repo)

	// 3. Wails Run
	err = wails.Run(&options.App{
		Title:     "Lag Monitor Pro",
		Width:     1024,
		Height:    768,
		Assets:    assets,
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
