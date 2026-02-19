package config

import (
	"encoding/json"
	"lag-monitor/internal/domain"
	"os"
	"path/filepath"
	"sync"
)

func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".config", "lagmon")
	os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "settings.json")
}

// Estrutura do Diagrama (Visualização)
type DiagramNode struct {
	Name string `json:"name"`
	IP   string `json:"ip"` // O IP é usado para vincular ao Ping
}

type NetworkDiagramConfig struct {
	Local    DiagramNode `json:"local"`
	Gateway  DiagramNode `json:"gateway"`
	Internet DiagramNode `json:"internet"`
}

// Estrutura Principal do Arquivo
type AppConfig struct {
	RetentionDays  int                  `json:"retention_days"`
	NetworkDiagram NetworkDiagramConfig `json:"network_diagram"`
	Targets        []domain.Host        `json:"targets"`
}

type ConfigManager struct {
	mu   sync.Mutex
	Data AppConfig
}

func NewConfigManager() *ConfigManager {
	return &ConfigManager{}
}

// Load lê o arquivo ou cria o padrão se não existir
func (c *ConfigManager) Load() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	path := GetConfigPath()
	file, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		// Cria Defaults
		c.Data = AppConfig{
			RetentionDays: 7,
			NetworkDiagram: NetworkDiagramConfig{
				Local:    DiagramNode{Name: "You (Local)", IP: "127.0.0.1"},
				Gateway:  DiagramNode{Name: "Gateway", IP: "192.168.1.1"},
				Internet: DiagramNode{Name: "Internet", IP: "8.8.8.8"},
			},
			Targets: []domain.Host{
				{ID: "gateway", Name: "Gateway", IP: "192.168.1.1", IsGW: true, Active: true},
				{ID: "google", Name: "Google DNS", IP: "8.8.8.8", IsGW: false, Active: true},
			},
		}
		return c.Save() // Cria o arquivo físico
	} else if err != nil {
		return err
	}

	return json.Unmarshal(file, &c.Data)
}

func (c *ConfigManager) Save() error {
	// A função Save é chamada internamente já com Lock, ou externamente.
	// Se for chamar externamente direto, descomentar o lock.
	// c.mu.Lock()
	// defer c.mu.Unlock()

	data, err := json.MarshalIndent(c.Data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(GetConfigPath(), data, 0644)
}

// Helpers para atualizar alvos e salvar
func (c *ConfigManager) AddTarget(h domain.Host) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Verifica duplicidade
	for _, t := range c.Data.Targets {
		if t.ID == h.ID {
			return nil
		}
	}
	c.Data.Targets = append(c.Data.Targets, h)
	return c.Save()
}

func (c *ConfigManager) RemoveTarget(id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	newTargets := []domain.Host{}
	for _, t := range c.Data.Targets {
		if t.ID != id {
			newTargets = append(newTargets, t)
		}
	}
	c.Data.Targets = newTargets
	return c.Save()
}

func (c *ConfigManager) UpdateTargetStatus(id string, active bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, t := range c.Data.Targets {
		if t.ID == id {
			c.Data.Targets[i].Active = active
			break
		}
	}
	return c.Save()
}

func (c *ConfigManager) UpdateConfig(newCfg AppConfig) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Data = newCfg
}

func (c *ConfigManager) UpdateTargetDiagramVisibility(id string, show bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, t := range c.Data.Targets {
		if t.ID == id {
			c.Data.Targets[i].ShowInDiagram = show
			break
		}
	}
	return c.Save()
}
