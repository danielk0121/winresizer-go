# 버그: IntelliJ에서 실행 시 "기본 세팅값으로 초기화" 실패

## 현상
IntelliJ에서 `main.go`를 직접 실행한 상태에서 웹 UI의 "기본 세팅값으로 초기화" 버튼을 누르면 에러가 표시되고 초기화되지 않는다.

## 원인
`resolveDefaultConfigPath()`에서 `default-config.json`을 두 가지 경로로 탐색한다.

1. **실행 파일 기준**: IntelliJ가 빌드한 바이너리는 `/private/var/folders/.../...` 같은 임시 경로에 생성되며, 해당 경로에 `config/` 폴더가 없다.
2. **CWD 기준 폴백** (`config/default-config.json`): IntelliJ의 기본 실행 CWD는 프로젝트 루트(`/Users/user/ws/winresizer-go`)이므로, `config/default-config.json`을 찾지 못한다. 실제 파일은 `app/config/default-config.json`에 있다.

두 경로 모두 실패하면 `LoadDefaultConfig`가 에러를 반환하고 `/api/config/reset`이 500 에러를 응답한다.

## 해결 방안
`runtime.Caller(0)`으로 소스 파일의 절대 경로를 얻어, 소스 기준 상대 경로(`../config/default-config.json`)를 추가 탐색 후보로 삽입한다.

```go
_, srcFile, _, ok := runtime.Caller(0)
if ok {
    // srcFile: .../app/core/config_manager.go
    candidate := filepath.Join(filepath.Dir(srcFile), "..", "config", "default-config.json")
    if _, err := os.Stat(candidate); err == nil {
        return candidate
    }
}
```

탐색 순서:
1. 실행 파일 기준 (배포 번들)
2. 소스 파일 기준 (IntelliJ 등 IDE 실행)
3. CWD 기준 폴백 (`go run ./...`)

## 수정 파일
- `app/core/config_manager.go`: `resolveDefaultConfigPath`에 소스 파일 기준 경로 탐색 추가
