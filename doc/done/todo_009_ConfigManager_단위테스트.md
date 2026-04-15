# TODO-009: ConfigManager 단위 테스트 작성

## 우선순위
4순위 (테스트 보완)

## 목표
`app/core/config_manager.go`에 대한 단위 테스트 작성.
(`spec_테스트계획.md`에 계획만 있고 미작성 상태)

## 테스트 파일 위치
`app/core/config_manager_test.go`

## 테스트 케이스 목록
| # | 테스트명 | 시나리오 |
|---|---------|---------|
| 1 | TestLoadDefaultConfig | default-config.json 정상 로드 |
| 2 | TestLoadConfig_NoUserConfig | config.json 없을 때 기본값으로 초기화 |
| 3 | TestLoadConfig_WithUserConfig | config.json 있을 때 사용자 설정 로드 |
| 4 | TestSaveConfig | 설정 저장 후 파일 내용 검증 |
| 5 | TestSaveRuntimeInfo | 런타임 정보 저장 검증 |
| 6 | TestInvalidateCache | 캐시 무효화 후 재로드 검증 |
| 7 | TestLoadConfig_CorruptFile | 손상된 config.json → 기본값 폴백 |

## 구현 방법
- 임시 디렉토리(`t.TempDir()`)를 `configFilePath`로 오버라이드하여 실제 파일 생성 없이 테스트
- `configFilePath` 변수를 테스트에서 주입 가능하도록 수정 필요

## 검증 방법
```bash
go test ./app/core/ -run TestLoadConfig -v
go test ./app/core/ -run TestSave -v
```

## 참고 문서
- `doc/done/spec_테스트계획.md`

## 작업 결과

**상태**: 완료

**생성 파일**: `app/core/config_manager_test.go`
**수정 파일**: `app/core/config_manager.go`

**테스트 구현**:
`setupConfigTest()` 헬퍼로 `configFilePath`와 `defaultConfigPathFn`을 임시 경로로 교체하여 실제 사용자 설정 파일에 영향 없이 독립 테스트 실행.

| 테스트 | 결과 |
|--------|------|
| TestLoadDefaultConfig | PASS |
| TestLoadConfig_NoUserConfig | PASS |
| TestLoadConfig_WithUserConfig | PASS |
| TestSaveConfig | PASS |
| TestSaveRuntimeInfo | PASS |
| TestInvalidateCache | PASS |
| TestLoadConfig_CorruptFile | PASS |

**버그 수정 (부산물)**:
`LoadConfig()`에서 `cacheMu.Lock()` 보유 중 `SaveConfig()` 호출 시 `SaveConfig()`도 `cacheMu.Lock()`을 시도해 데드락 발생. `saveConfigFile()` 내부 함수를 분리하여 해결.

**커밋**: `886a477`
