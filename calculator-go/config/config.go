package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type ServerConfig struct {
	Port string `yaml:"port"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
}

// LoadConfig 加载配置文件，如果文件不存在则使用默认配置
func LoadConfig() (*Config, error) {
	// 设置默认配置
	defaultConfig := &Config{
		Server: ServerConfig{
			Port: "8080", // 默认端口
		},
	}

	// 尝试从多个位置查找配置文件
	configPaths := []string{
		"config.yaml",                          // 当前目录
		filepath.Join("config", "config.yaml"), // config目录下
	}

	var configFile string
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			configFile = path
			break
		}
	}

	// 如果找不到配置文件，返回默认配置
	if configFile == "" {
		return defaultConfig, nil
	}

	// 读取配置文件
	file, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	// 解析YAML
	loadedConfig := defaultConfig // 使用默认值作为基础
	if err := yaml.Unmarshal(file, loadedConfig); err != nil {
		return nil, err
	}

	return loadedConfig, nil
}

// MustLoadConfig 加载配置，如果出错则panic（适合在程序启动时使用）
func MustLoadConfig() *Config {
	cfg, err := LoadConfig()
	if err != nil {
		panic("加载配置文件失败: " + err.Error())
	}
	return cfg
}
