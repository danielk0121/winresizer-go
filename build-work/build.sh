#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
APP_DIR="${SCRIPT_DIR}/../app"

APP_NAME="WinResizer"
BUNDLE_DIR="${SCRIPT_DIR}/dist/${APP_NAME}.app"
CONTENTS_DIR="${BUNDLE_DIR}/Contents"
MACOS_DIR="${CONTENTS_DIR}/MacOS"
RESOURCES_DIR="${CONTENTS_DIR}/Resources"

echo "==> Cleaning dist/"
rm -rf "${SCRIPT_DIR}/dist/"
mkdir -p "${MACOS_DIR}" "${RESOURCES_DIR}"

echo "==> Building Universal Binary (arm64 + amd64)"
cd "${APP_DIR}"
GOARCH=arm64 CGO_ENABLED=1 go build -o "${SCRIPT_DIR}/dist/winresizer_arm64" .
GOARCH=amd64 CGO_ENABLED=1 go build -o "${SCRIPT_DIR}/dist/winresizer_amd64" .
lipo -create "${SCRIPT_DIR}/dist/winresizer_arm64" "${SCRIPT_DIR}/dist/winresizer_amd64" -output "${MACOS_DIR}/${APP_NAME}"
rm "${SCRIPT_DIR}/dist/winresizer_arm64" "${SCRIPT_DIR}/dist/winresizer_amd64"

echo "==> Copying resources"
cp "${SCRIPT_DIR}/Info.plist" "${CONTENTS_DIR}/Info.plist"
if [ -f "${APP_DIR}/ui/icon.icns" ]; then
    cp "${APP_DIR}/ui/icon.icns" "${RESOURCES_DIR}/icon.icns"
fi

# 빌드 시간 마커 파일 생성
BUILD_TIME=$(date +"%Y%m%d-%H%M")
touch "${CONTENTS_DIR}/buildtime-${BUILD_TIME}.txt"

echo "==> Code signing (ad-hoc)"
codesign --force --deep --sign - \
    --entitlements "${SCRIPT_DIR}/entitlements.plist" \
    "${BUNDLE_DIR}"

echo "==> Creating DMG"
# 드래그 설치용 임시 폴더 구성
mkdir -p "${SCRIPT_DIR}/dist/dmg"
cp -r "${BUNDLE_DIR}" "${SCRIPT_DIR}/dist/dmg/"
ln -sf /Applications "${SCRIPT_DIR}/dist/dmg/Applications"

hdiutil create -volname "${APP_NAME}" \
    -srcfolder "${SCRIPT_DIR}/dist/dmg" \
    -ov -format UDZO \
    "${SCRIPT_DIR}/dist/${APP_NAME}.dmg"

rm -rf "${SCRIPT_DIR}/dist/dmg"

echo "==> Done: dist/${APP_NAME}.dmg"
