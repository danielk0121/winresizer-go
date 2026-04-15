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
