#!/bin/bash
# GitHub 릴리즈 배포 스크립트
# 사용법:
#   ./release.sh v1.0.0
#   ./release.sh v1.0.0 "릴리즈 노트 내용"
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
DMG_PATH="${SCRIPT_DIR}/dist/WinResizer.dmg"

# ── 인자 검증 ────────────────────────────────────────────────────
TAG="${1}"
NOTES="${2:-}"

if [ -z "${TAG}" ]; then
    echo "사용법: $0 <태그> [릴리즈노트]"
    echo "예시:  $0 v1.0.0"
    echo "예시:  $0 v1.0.0 \"버그 수정 및 기능 개선\""
    exit 1
fi

# ── DMG 존재 확인 ────────────────────────────────────────────────
if [ ! -f "${DMG_PATH}" ]; then
    echo "오류: DMG 파일이 없습니다 — ${DMG_PATH}"
    echo "먼저 build.sh 를 실행하세요."
    exit 1
fi

# ── gh 설치 확인 ─────────────────────────────────────────────────
if ! command -v gh &>/dev/null; then
    echo "오류: gh (GitHub CLI) 가 설치되어 있지 않습니다."
    echo "설치: brew install gh"
    exit 1
fi

echo "==> 태그: ${TAG}"
echo "==> DMG:  ${DMG_PATH}"

# ── 릴리즈 생성 및 DMG 업로드 ────────────────────────────────────
if [ -n "${NOTES}" ]; then
    gh release create "${TAG}" "${DMG_PATH}" \
        --repo danielk0121/winresizer-go \
        --title "${TAG}" \
        --notes "${NOTES}"
else
    gh release create "${TAG}" "${DMG_PATH}" \
        --repo danielk0121/winresizer-go \
        --title "${TAG}" \
        --generate-notes
fi

echo "==> 릴리즈 완료: https://github.com/danielk0121/winresizer-go/releases/tag/${TAG}"
