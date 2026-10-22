package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/toolkits/file"
)

// LogConfig 日志配置
type LogConfig struct {
	Level string `json:"level"`
}

// JwtConfig JWT配置
type JwtConfig struct {
	TokenName string `json:"token_name"` // "Authentication"
	Issuer    string `json:"issuer"`     // "db"
	SecretKey string `json:"secret_key"` // "362b36fd7d514c519333234da152eaff"
	TTL       string `json:"ttl"`        // Go duration, empty means 24h
}

// SecretConfig 安全配置
type SecretConfig struct {
	Jwt    *JwtConfig `json:"jwt"`
	Crypto string     `json:"crypto"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Addr           string `json:"addr"`
	MaxIdle        int    `json:"max_idle"`
	MaxConnections int    `json:"max_connections"`
	WaitTimeout    int    `json:"wait_timeout"`
	// Driver selects the metadata engine: "mysql" or "postgres".
	// Empty means auto-detect from Addr (postgres:// URL => postgres).
	Driver string `json:"driver"`
}

// GlobalConfig 配置
type GlobalConfig struct {
	Log      *LogConfig      `json:"log"`
	Cert     string          `json:"cert"`
	Key      string          `json:"key"`
	Database *DatabaseConfig `json:"database"`
	Listen   string          `json:"listen"`
	Secret   *SecretConfig   `json:"secret"`
}

var (
	// ConfigFile 实际读到的配置文件，没有则为空。
	ConfigFile string
	// SavePath 配置写入位置：向导完成后把完整配置写到这里。
	SavePath   string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

// Config 返回当前的配置
func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

// ParseConfig 从配置文件读取配置，反序列化成配置对象。
//
// 查找顺序：-c 指定的文件 → 运行目录的 cfg.json → /etc/halo/cfg.json。
// 都没有就用默认值启动，由 /welcome 向导写出配置文件。
func ParseConfig(cfg string) {
	configLock.Lock()
	ConfigFile = ""
	SavePath = ""
	configLock.Unlock()

	candidates := []string{}
	if trimmed := strings.TrimSpace(cfg); trimmed != "" {
		candidates = append(candidates, trimmed)
	} else {
		if cwd, err := os.Getwd(); err == nil {
			candidates = append(candidates, filepath.Join(cwd, "cfg.json"))
		}
		candidates = append(candidates, "/etc/halo/cfg.json")
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			ConfigFile = candidate
			log.Debugf("[D] Loading config from: %s", candidate)
			break
		}
	}
	// 写入位置：-c 指定的文件优先，否则写运行目录，再不行写 /etc/halo。
	SavePath = candidates[0]
	if len(candidates) > 1 && !writable(filepath.Dir(SavePath)) {
		for _, candidate := range candidates[1:] {
			if writable(filepath.Dir(candidate)) {
				SavePath = candidate
				break
			}
		}
	}

	if ConfigFile == "" {
		// 没有配置文件：先用默认值跑起来，向导完成后再落盘。
		log.Infof("[I] 未找到配置文件，使用默认配置，初始化向导会把配置写入 %s", SavePath)
		configLock.Lock()
		config = &GlobalConfig{}
		configLock.Unlock()
		applyDefaults()
		return
	}

	content, err := file.ToTrimString(ConfigFile)
	if err != nil {
		log.Fatalf("[F] 读取配置文件 \"%s\" 错误: %s", ConfigFile, err.Error())
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(content), &c)
	if err != nil {
		log.Fatalf("[F] 解析配置文件 \"%s\" 错误: %s", ConfigFile, err.Error())
	}

	configLock.Lock()
	config = &c
	configLock.Unlock()

	applyDefaults()

	log.Debugf("[D] 读取配置文件 \"%s\" 成功。", ConfigFile)
}

// writable reports whether the directory accepts new files.
func writable(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return false
	}
	probe, err := os.CreateTemp(dir, ".halo-probe-*")
	if err != nil {
		return false
	}
	name := probe.Name()
	_ = probe.Close()
	_ = os.Remove(name)
	return true
}

// applyDefaults fills in the values the server always needs.
func applyDefaults() {
	configLock.Lock()
	defer configLock.Unlock()
	if config.Log == nil {
		config.Log = &LogConfig{}
	}
	if config.Log.Level == "" {
		config.Log.Level = "info"
	}
	if config.Cert == "" {
		config.Cert = "certs/cert.pem"
	}
	if config.Key == "" {
		config.Key = "certs/key.pem"
	}
	if config.Listen == "" {
		config.Listen = ":8443"
	}
}

// Configured reports whether the metadata database and the secrets are set.
// A missing part means the server has to run in setup mode.
func Configured() bool {
	cfg := Config()
	if cfg == nil {
		return false
	}
	if cfg.Database == nil || strings.TrimSpace(cfg.Database.Addr) == "" {
		return false
	}
	if cfg.Secret == nil || cfg.Secret.Jwt == nil {
		return false
	}
	if strings.TrimSpace(cfg.Secret.Jwt.SecretKey) == "" {
		return false
	}
	return strings.TrimSpace(cfg.Secret.Crypto) != ""
}

// Save applies the mutation to the current configuration and writes it back
// to the config file. The file is written through a temporary file so a
// crash cannot leave a half-written config behind.
func Save(mutate func(*GlobalConfig)) error {
	if mutate == nil {
		return fmt.Errorf("config: 没有要保存的变更")
	}
	if strings.TrimSpace(SavePath) == "" {
		return fmt.Errorf("config: 配置文件路径未知")
	}

	// 复制一份再改：指针字段要单独复制，否则会改到正在使用的配置。
	configLock.RLock()
	current := GlobalConfig{}
	if config != nil {
		current = *config
	}
	if current.Database != nil {
		database := *current.Database
		current.Database = &database
	}
	if current.Secret != nil {
		secret := *current.Secret
		if current.Secret.Jwt != nil {
			jwt := *current.Secret.Jwt
			secret.Jwt = &jwt
		}
		current.Secret = &secret
	}
	configLock.RUnlock()

	mutate(&current)

	body, err := json.MarshalIndent(&current, "", "  ")
	if err != nil {
		return fmt.Errorf("config: 序列化配置失败: %w", err)
	}
	body = append(body, '\n')

	dir := filepath.Dir(SavePath)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("config: 创建配置目录失败: %w", err)
	}
	tmp, err := os.CreateTemp(dir, ".cfg-*.json")
	if err != nil {
		return fmt.Errorf("config: 无法在 %s 写入配置文件，该目录不可写，请用 -c 指定可写路径: %w", dir, err)
	}
	tmpName := tmp.Name()
	defer func() { _ = os.Remove(tmpName) }()

	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("config: 设置配置文件权限失败: %w", err)
	}
	if _, err := tmp.Write(body); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("config: 写入配置文件失败: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("config: 关闭临时配置文件失败: %w", err)
	}
	if err := os.Rename(tmpName, SavePath); err != nil {
		return fmt.Errorf("config: 替换配置文件 %s 失败: %w", SavePath, err)
	}
	ConfigFile = SavePath

	configLock.Lock()
	config = &current
	configLock.Unlock()
	applyDefaults()
	return nil
}
