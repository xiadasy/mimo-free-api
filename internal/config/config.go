
package config

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
)

type Config struct {
	Port         string    `json:"port"`
	APIKey       string    `json:"api_key"` // 网关自身的认证 Key
	DefaultModel string    `json:"default_model"`
	Accounts     []Account `json:"accounts"`
}

// Account 网页端账号配置
type Account struct {
	ID           string `json:"id"`
	ServiceToken string `json:"service_token"`
	UserID       string `json:"user_id"`
	Ph           string `json:"ph"`
	Active       bool     `json:"active"`
}

var (
	cfg  *Config
	mu   sync.RWMutex
	path string
)

func Load(p string) (*Config, error) {
	path = p
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg = &Config{Port: "8080", APIKey: "sk-mimo", DefaultModel: "mimo-v2.5-pro"}
			applyEnvOverrides()
			return cfg, Save()
		}
		return nil, err
	}
	cfg = &Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	applyEnvOverrides()
	return cfg, nil
}

func applyEnvOverrides() {
	if cfg == nil {
		return
	}
	if v := envValue("PORT"); v != "" {
		cfg.Port = v
	}
	if v := envValue("MIMO_API_KEY"); v != "" {
		cfg.APIKey = v
	} else if v := envValue("API_KEY"); v != "" {
		cfg.APIKey = v
	}
	if v := envValue("MIMO_DEFAULT_MODEL"); v != "" {
		cfg.DefaultModel = v
	}

	serviceToken := envValue("MIMO_SERVICE_TOKEN")
	userID := envValue("MIMO_USER_ID")
	ph := envValue("MIMO_PH")
	if serviceToken != "" && userID != "" && ph != "" {
		id := envValue("MIMO_ACCOUNT_ID")
		if id == "" {
			id = "env-account-1"
		}
		account := Account{ID: id, ServiceToken: serviceToken, UserID: userID, Ph: ph, Active: true}
		replaced := false
		for i := range cfg.Accounts {
			if cfg.Accounts[i].ID == id {
				cfg.Accounts[i] = account
				replaced = true
				break
			}
		}
		if !replaced {
			cfg.Accounts = append(cfg.Accounts, account)
		}
	}
}

func envValue(key string) string {
	return strings.Trim(strings.TrimSpace(os.Getenv(key)), "\"")
}

func Get() Config {
	mu.RLock()
	defer mu.RUnlock()
	return *cfg
}

func Save() error {
	mu.RLock()
	data, err := json.MarshalIndent(cfg, "", "  ")
	mu.RUnlock()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func Update(fn func(*Config)) {
	mu.Lock()
	fn(cfg)
	mu.Unlock()
}
