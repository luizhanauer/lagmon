package main

import (
	"embed"
	"lag-monitor/internal/config" // Importe o novo pacote config
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
	// 1. Infra - Banco de Dados
	repo, err := database.NewSQLiteRepo("./lagmonitor.db")
	if err != nil {
		log.Fatal(err)
	}

	// 2. Infra - Configuração (NOVO)
	// Instanciamos o gerenciador de configuração aqui
	cfg := config.NewConfigManager()

	// 3. Infra - Rede
	pinger := network.NewPinger()

	// 4. Serviço
	// Usamos uma variável declarada antes para o closure do emitter capturar o contexto do App
	var app *App

	emitter := func(event string, data interface{}) {
		// Verifica se o app e o contexto já foram inicializados
		if app != nil && app.ctx != nil {
			runtime.EventsEmit(app.ctx, event, data)
		}
	}

	service := usecase.NewMonitorService(repo, pinger, emitter)

	// 5. Inicialização do App (ATUALIZADO)
	// Agora passamos o repo e o cfg para dentro do App
	app = NewApp(service, repo, cfg)

	// 6. Wails Run
	err = wails.Run(&options.App{
		Title:     "LAGMON",
		Width:     1024, // Ajustei para um tamanho inicial mais confortável para o Dashboard
		Height:    768,
		Assets:    assets,
		OnStartup: app.startup, // O app.startup agora vai chamar cfg.Load()
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
