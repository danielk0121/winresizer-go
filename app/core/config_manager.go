package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
	"winresizer/utils"
)

// 설정 파일 경로
var configFilePath = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "WinResizer", "config.json")

// ShortcutConfig는 단축키 하나의 설정입니다.
type ShortcutConfig struct {
	Display   string `json:"display"`
	Mode      string `json:"mode"`
	Keycode   uint32 `json:"keycode"`
	Modifiers uint32 `json:"modifiers"`
}

// Settings는 앱 일반 설정입니다.
type Settings struct {
	Gap          int               `json:"gap"`
	LoginLaunch  bool              `json:"login_launch"`
	IgnoreApps   []string          `json:"ignore_apps"`
	AutoLayouts  map[string]string `json:"auto_layouts"`
}

// RuntimeInfo는 앱 실행 중 기록되는 런타임 정보입니다.
type RuntimeInfo struct {
	Port      int    `json:"port"`
	PID       int    `json:"pid"`
	StartTime string `json:"start_time"`
}

// Config는 전체 설정 구조체입니다.
type Config struct {
	Runtime   RuntimeInfo               `json:"runtime,omitempty"`
	Settings  Settings                  `json:"settings"`
	Shortcuts map[string]ShortcutConfig `json:"shortcuts"`
}

var (
	configCache *Config
	cacheMu     sync.RWMutex
)

// defaultConfigPath는 번들 내 기본 설정 파일 경로를 반환합니다.
func defaultConfigPath() string {
	// 실행 파일 기준 상대 경로로 탐색
	exe, err := os.Executable()
	if err == nil {
		candidate := filepath.Join(filepath.Dir(exe), "config", "default-config.json")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	// 개발 환경: 소스 기준 경로
	return filepath.Join("config", "default-config.json")
}

// LoadDefaultConfig는 default-config.json을 읽어 반환합니다.
func LoadDefaultConfig() (*Config, error) {
	data, err := os.ReadFile(defaultConfigPath())
	if err != nil {
		return nil, fmt.Errorf("기본 설정 파일 로드 실패: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("기본 설정 파일 파싱 실패: %w", err)
	}
	return &cfg, nil
}

// ensureConfigDir은 설정 파일 디렉토리를 생성합니다.
func ensureConfigDir() error {
	return os.MkdirAll(filepath.Dir(configFilePath), 0755)
}

// LoadConfig는 설정을 불러옵니다.
// config.json이 있으면 로드, 없으면 default-config.json으로 초기화합니다.
func LoadConfig() (*Config, error) {
	cacheMu.RLock()
	if configCache != nil {
		defer cacheMu.RUnlock()
		return configCache, nil
	}
	cacheMu.RUnlock()

	cacheMu.Lock()
	defer cacheMu.Unlock()

	// 사용자 설정 파일이 있으면 로드
	if _, err := os.Stat(configFilePath); err == nil {
		data, err := os.ReadFile(configFilePath)
		if err == nil {
			var cfg Config
			if err := json.Unmarshal(data, &cfg); err == nil {
				configCache = &cfg
				return configCache, nil
			}
		}
		utils.Log.Errorf("사용자 설정 파일 로드 실패, 기본값으로 초기화합니다.")
	}

	// 없으면 기본값으로 초기화
	utils.Log.Infof("사용자 설정 파일이 없어 기본값으로 초기화합니다.")
	cfg, err := LoadDefaultConfig()
	if err != nil {
		return nil, err
	}
	if err := SaveConfig(cfg); err != nil {
		utils.Log.Errorf("초기 설정 저장 실패: %v", err)
	}
	configCache = cfg
	return configCache, nil
}

// GetConfig는 캐시된 설정을 반환합니다.
func GetConfig() (*Config, error) {
	return LoadConfig()
}

// SaveConfig는 설정을 config.json에 저장합니다.
func SaveConfig(cfg *Config) error {
	if err := ensureConfigDir(); err != nil {
		return fmt.Errorf("설정 디렉토리 생성 실패: %w", err)
	}
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return fmt.Errorf("설정 직렬화 실패: %w", err)
	}
	if err := os.WriteFile(configFilePath, data, 0644); err != nil {
		return fmt.Errorf("설정 파일 저장 실패: %w", err)
	}
	// 캐시 갱신
	cacheMu.Lock()
	configCache = cfg
	cacheMu.Unlock()
	utils.Log.Debugf("설정 파일 저장 완료: %s", configFilePath)
	return nil
}

// InvalidateCache는 설정 캐시를 무효화합니다. 단축키 리스너 재시작 시 호출합니다.
func InvalidateCache() {
	cacheMu.Lock()
	configCache = nil
	cacheMu.Unlock()
}

// SaveRuntimeInfo는 런타임 정보(port, pid, start_time)를 설정 파일에 기록합니다.
func SaveRuntimeInfo(port int) error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}
	cfg.Runtime = RuntimeInfo{
		Port:      port,
		PID:       os.Getpid(),
		StartTime: time.Now().Format(time.RFC3339),
	}
	return SaveConfig(cfg)
}
