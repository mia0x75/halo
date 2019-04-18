package g

import (
	"encoding/json"
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

	for _, config := range configs {
		if _, err = os.Stat(config); err == nil {
			ConfigFile = config
			log.Debugf("[D] Loading config from: %s", config)
			break
		}
	}
	if err != nil {
		log.Fatalf("[F] 读取配置文件错误。")
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
	defer configLock.Unlock()

	config = &c
	if config.Cert == "" {
		config.Cert = "cert.pem"
	}
	if config.Key == "" {
		config.Key = "key.pem"
	}

	log.Debugf("[D] 读取配置文件 \"%s\" 成功。", ConfigFile)
}
