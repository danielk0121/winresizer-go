# TODO-008: Universal Binary 빌드 (arm64 + amd64)

## 우선순위
3순위 (배포 준비)

## 목표
Apple Silicon(arm64)과 Intel(amd64) Mac 모두에서 동작하는 단일 바이너리 배포.

## 구현 방법
```bash
# arm64 빌드
GOARCH=arm64 CGO_ENABLED=1 go build -o dist/winresizer_arm64 ./app

# amd64 빌드 (크로스 컴파일 — CGo 크로스 컴파일러 필요)
GOARCH=amd64 CGO_ENABLED=1 \
  CC=x86_64-apple-darwin-clang \
  go build -o dist/winresizer_amd64 ./app

# lipo로 통합
lipo -create dist/winresizer_arm64 dist/winresizer_amd64 \
  -output WinResizer.app/Contents/MacOS/winresizer
```

### CGo 크로스 컴파일 주의사항
- macOS에서 `x86_64-apple-darwin` 크로스 컴파일러 설치 필요
- `osxcross` 툴체인 또는 Xcode 명령줄 도구 사용

## 의존성
- TODO-006 (앱 번들 빌드) 완료 후 진행

## 검증 방법
```bash
lipo -info WinResizer.app/Contents/MacOS/winresizer
# 출력: Architectures in the fat file: arm64 x86_64
```
