package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// setupConfigTest는 테스트용 임시 경로를 설정하고 캐시를 초기화합니다.
// 반환된 cleanup 함수를 defer로 호출해야 합니다.
func setupConfigTest(t *testing.T) (tmpDir string, cleanup func()) {
	t.Helper()
	tmpDir = t.TempDir()
	origPath := configFilePath
	origDefaultFn := defaultConfigPathFn

	configFilePath = filepath.Join(tmpDir, "config.json")
	// default-config.json 경로를 소스 기준 절대 경로로 설정 (core 패키지 기준 상위)
	defaultConfigPathFn = func() string {
		return filepath.Join("..", "config", "default-config.json")
	}
	InvalidateCache()
	return tmpDir, func() {
		configFilePath = origPath
		defaultConfigPathFn = origDefaultFn
		InvalidateCache()
	}
}

func TestLoadDefaultConfig(t *testing.T) {
	_, cleanup := setupConfigTest(t)
	defer cleanup()

	cfg, err := LoadDefaultConfig()
	if err != nil {
		t.Fatalf("LoadDefaultConfig 실패: %v", err)
	}
	if len(cfg.Shortcuts) == 0 {
		t.Error("Shortcuts가 비어 있습니다")
	}
	// 기본 단축키 Left가 존재하는지 확인
	if _, ok := cfg.Shortcuts["Left"]; !ok {
		t.Error("기본 단축키 'Left'가 없습니다")
	}
}

func TestLoadConfig_NoUserConfig(t *testing.T) {
	_, cleanup := setupConfigTest(t)
	defer cleanup()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig 실패: %v", err)
	}
	if cfg == nil {
		t.Fatal("cfg가 nil입니다")
	}
	// 기본값으로 초기화 → Shortcuts 존재 확인
	if len(cfg.Shortcuts) == 0 {
		t.Error("기본값 초기화 후 Shortcuts가 비어 있습니다")
	}
	// config.json이 생성됐는지 확인
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		t.Error("config.json이 생성되지 않았습니다")
	}
}

func TestLoadConfig_WithUserConfig(t *testing.T) {
	tmpDir, cleanup := setupConfigTest(t)
	defer cleanup()

	// 사용자 config.json을 미리 작성
	userCfg := Config{
		Settings: Settings{Gap: 8},
		Shortcuts: map[string]ShortcutConfig{
			"Left": {Display: "custom", Mode: "left_half", Keycode: 100, Modifiers: 9999},
		},
	}
	data, _ := json.MarshalIndent(userCfg, "", "    ")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configFilePath, data, 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig 실패: %v", err)
	}
	if cfg.Settings.Gap != 8 {
		t.Errorf("사용자 설정 Gap 로드 실패: got %d, want 8", cfg.Settings.Gap)
	}
	if cfg.Shortcuts["Left"].Keycode != 100 {
		t.Errorf("사용자 설정 Keycode 로드 실패: got %d, want 100", cfg.Shortcuts["Left"].Keycode)
	}
}

func TestSaveConfig(t *testing.T) {
	_, cleanup := setupConfigTest(t)
	defer cleanup()

	cfg := &Config{
		Settings: Settings{Gap: 4, LoginLaunch: true},
		Shortcuts: map[string]ShortcutConfig{
			"Right": {Display: "test", Mode: "right_half", Keycode: 124, Modifiers: 6400},
		},
	}
	if err := SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig 실패: %v", err)
	}

	// 파일 내용 검증
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		t.Fatalf("저장된 파일 읽기 실패: %v", err)
	}
	var loaded Config
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("저장된 파일 파싱 실패: %v", err)
	}
	if loaded.Settings.Gap != 4 {
		t.Errorf("저장된 Gap 불일치: got %d, want 4", loaded.Settings.Gap)
	}
	if loaded.Shortcuts["Right"].Keycode != 124 {
		t.Errorf("저장된 Keycode 불일치: got %d, want 124", loaded.Shortcuts["Right"].Keycode)
	}
}

func TestSaveRuntimeInfo(t *testing.T) {
	_, cleanup := setupConfigTest(t)
	defer cleanup()

	if err := SaveRuntimeInfo(8080); err != nil {
		t.Fatalf("SaveRuntimeInfo 실패: %v", err)
	}

	// 캐시 무효화 후 다시 로드해서 확인
	InvalidateCache()
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig 실패: %v", err)
	}
	if cfg.Runtime.Port != 8080 {
		t.Errorf("런타임 포트 저장 실패: got %d, want 8080", cfg.Runtime.Port)
	}
	if cfg.Runtime.PID != os.Getpid() {
		t.Errorf("런타임 PID 저장 실패: got %d, want %d", cfg.Runtime.PID, os.Getpid())
	}
	if cfg.Runtime.StartTime == "" {
		t.Error("StartTime이 비어 있습니다")
	}
}

func TestInvalidateCache(t *testing.T) {
	_, cleanup := setupConfigTest(t)
	defer cleanup()

	// 1차 로드 → 캐시에 저장됨
	cfg1, err := LoadConfig()
	if err != nil {
		t.Fatalf("1차 LoadConfig 실패: %v", err)
	}
	origGap := cfg1.Settings.Gap

	// 파일만 직접 변경 (saveConfigFile 사용하여 캐시 우회)
	modified := *cfg1
	modified.Settings.Gap = 99
	_ = saveConfigFile(&modified)

	// 캐시 무효화 없이 로드 → 캐시된 원래 값 반환
	cfg2, _ := LoadConfig()
	if cfg2.Settings.Gap != origGap {
		t.Errorf("캐시 무효화 전에 파일 변경값이 보이면 안 됩니다: got %d, want %d", cfg2.Settings.Gap, origGap)
	}

	// 캐시 무효화 후 로드 → 파일에서 새 값 반환
	InvalidateCache()
	cfg3, err := LoadConfig()
	if err != nil {
		t.Fatalf("무효화 후 LoadConfig 실패: %v", err)
	}
	if cfg3.Settings.Gap != 99 {
		t.Errorf("캐시 무효화 후 새 값 로드 실패: got %d, want 99", cfg3.Settings.Gap)
	}
}

func TestLoadConfig_CorruptFile(t *testing.T) {
	tmpDir, cleanup := setupConfigTest(t)
	defer cleanup()

	// 손상된 JSON 파일 작성
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configFilePath, []byte("{invalid json!!!}"), 0644); err != nil {
		t.Fatal(err)
	}

	// 손상된 파일 → 기본값으로 폴백
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("손상된 파일에서 LoadConfig 실패: %v", err)
	}
	if cfg == nil {
		t.Fatal("폴백 cfg가 nil입니다")
	}
	// 기본값으로 초기화됐으므로 Shortcuts 존재 확인
	if len(cfg.Shortcuts) == 0 {
		t.Error("기본값 폴백 후 Shortcuts가 비어 있습니다")
	}
}
