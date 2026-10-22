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
	// ConfigFile 配置文件
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

// Config 返回当前的配置
func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

// ParseConfig 从配置文件读取配置，反序列化成配置对象
func ParseConfig(cfg string) {
	var configs []string
	var path, content string
	var err error

	path, err = os.Executable()
	if err != nil {
		log.Fatalf("[F] 错误信息: %s", err.Error())
	}
	baseDir := filepath.Dir(path)

	// 指定了配置文件优先读配置文件，未指定配置文件按如下顺序加载，先找到哪个加载哪个
	if strings.TrimSpace(cfg) == "" {
		configs = []string{
			"/etc/halo/cfg.json",
			filepath.Join(baseDir, "etc", "cfg.json"),
			filepath.Join(baseDir, "cfg.json"),
		}
	} else {
		configs = []string{cfg}
	}

	for _, candidate := range configs {
		if _, err = os.Stat(candidate); err == nil {
			ConfigFile = candidate
			log.Debugf("[D] Loading config from: %s", candidate)
			break
		}
	}
	if ConfigFile == "" {
		// 配置文件还不存在：首次安装走 /welcome 向导，写入时才真正落盘。
		ConfigFile = firstWritablePath(configs)
		log.Warnf("[W] 未找到配置文件，将新建 %s 并进入初始化向导。", ConfigFile)
		configLock.Lock()
		config = &GlobalConfig{}
		configLock.Unlock()
		applyDefaults()
		return
	}

	content, err = file.ToTrimString(ConfigFile)
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

// applyDefaults fills in the values the server always needs.
func applyDefaults() {
	configLock.Lock()
	defer configLock.Unlock()
	if config.Cert == "" {
		config.Cert = "cert.pem"
	}
	if config.Key == "" {
		config.Key = "key.pem"
	}
	if config.Listen == "" {
		config.Listen = ":8443"
	}
}

// firstWritablePath returns the first candidate whose directory exists and
// accepts new files, so the setup wizard can create the config file.
func firstWritablePath(candidates []string) string {
	for _, candidate := range candidates {
		dir := filepath.Dir(candidate)
		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() {
			continue
		}
		probe, err := os.CreateTemp(dir, ".halo-probe-*")
		if err != nil {
			continue
		}
		name := probe.Name()
		_ = probe.Close()
		_ = os.Remove(name)
		return candidate
	}
	return candidates[len(candidates)-1]
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
	if strings.TrimSpace(ConfigFile) == "" {
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

	dir := filepath.Dir(ConfigFile)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("config: 创建配置目录失败: %w", err)
	}
	tmp, err := os.CreateTemp(dir, ".cfg-*.json")
	if err != nil {
		return fmt.Errorf("config: 创建临时配置文件失败: %w", err)
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
	if err := os.Rename(tmpName, ConfigFile); err != nil {
		return fmt.Errorf("config: 替换配置文件失败: %w", err)
	}

	configLock.Lock()
	config = &current
	configLock.Unlock()
	applyDefaults()
	return nil
}
