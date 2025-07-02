package xcomp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type ConfigService struct {
	config      map[string]any
	envMap      map[string]string
	mu          sync.RWMutex
	viper       *viper.Viper
	envPrefix   string
	initialized bool
}

// ConfigOptions for advanced configuration
type ConfigOptions struct {
	EnvPrefix    string
	EnvSeparator string
	AutoReload   bool
}

func NewConfigService(configPaths ...string) *ConfigService {
	opts := ConfigOptions{
		EnvPrefix:    "",
		EnvSeparator: "__",
		AutoReload:   false,
	}

	cs := &ConfigService{
		config:    make(map[string]any),
		envMap:    make(map[string]string),
		viper:     viper.New(),
		envPrefix: opts.EnvPrefix,
	}

	// Load .env file
	godotenv.Load()

	// Load environment variables
	cs.loadEnvironmentVariables(opts)

	// Load config files
	for _, configPath := range configPaths {
		cs.loadConfigFile(configPath)
	}

	cs.initialized = true
	return cs
}

func (cs *ConfigService) loadEnvironmentVariables(opts ConfigOptions) {
	// Setup viper for environment variables
	if cs.envPrefix != "" {
		cs.viper.SetEnvPrefix(cs.envPrefix)
	}
	cs.viper.SetEnvKeyReplacer(strings.NewReplacer(".", opts.EnvSeparator))
	cs.viper.AutomaticEnv()

	// Load all environment variables into envMap for backward compatibility
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			cs.envMap[pair[0]] = pair[1]
		}
	}
}

func (cs *ConfigService) loadConfigFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	var fileConfig map[string]any
	ext := filepath.Ext(path)

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &fileConfig); err != nil {
			return fmt.Errorf("failed to parse YAML config %s: %w", path, err)
		}
	default:
		return fmt.Errorf("unsupported config file format: %s (only .yaml/.yml supported)", ext)
	}

	cs.mergeConfig(fileConfig)

	// Also load into viper for advanced env override support
	configBuffer, _ := json.Marshal(fileConfig)
	cs.viper.ReadConfig(bytes.NewBuffer(configBuffer))

	return nil
}

func (cs *ConfigService) mergeConfig(newConfig map[string]any) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	for key, value := range newConfig {
		cs.config[key] = value
	}
}

func (cs *ConfigService) Get(key string) any {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	// Try viper first (supports env overrides with prefixes)
	if cs.initialized && cs.viper.IsSet(key) {
		return cs.viper.Get(key)
	}

	// Fallback to direct env lookup for backward compatibility
	if envValue, exists := cs.envMap[strings.ToUpper(key)]; exists {
		return envValue
	}

	if envValue, exists := cs.envMap[key]; exists {
		return envValue
	}

	return cs.getNestedValue(key)
}

func (cs *ConfigService) GetString(key string, defaultValue ...string) string {
	value := cs.Get(key)
	if value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}

	if str, ok := value.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", value)
}

func (cs *ConfigService) GetInt(key string, defaultValue ...int) int {
	value := cs.Get(key)
	if value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	switch v := value.(type) {
	case int:
		return v
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	case float64:
		return int(v)
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

func (cs *ConfigService) GetBool(key string, defaultValue ...bool) bool {
	value := cs.Get(key)
	if value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}

	switch v := value.(type) {
	case bool:
		return v
	case string:
		return strings.ToLower(v) == "true" || v == "1"
	case int:
		return v != 0
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return false
}

func (cs *ConfigService) getNestedValue(key string) any {
	keys := strings.Split(key, ".")
	current := cs.config

	for i, k := range keys {
		if i == len(keys)-1 {
			return current[k]
		}

		if next, ok := current[k].(map[string]any); ok {
			current = next
		} else {
			return nil
		}
	}

	return nil
}

func (cs *ConfigService) GetAll() map[string]any {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	result := make(map[string]any)
	for k, v := range cs.config {
		result[k] = v
	}
	for k, v := range cs.envMap {
		result[k] = v
	}
	return result
}
