#!/bin/bash
set -e

APP_NAME="WinResizer"
BUNDLE_DIR="dist/${APP_NAME}.app"
CONTENTS_DIR="${BUNDLE_DIR}/Contents"
MACOS_DIR="${CONTENTS_DIR}/MacOS"
RESOURCES_DIR="${CONTENTS_DIR}/Resources"

echo "==> Cleaning dist/"
rm -rf dist/
mkdir -p "${MACOS_DIR}" "${RESOURCES_DIR}"

echo "==> Building binary (현재 아키텍처: $(uname -m))"
go build -o "${MACOS_DIR}/${APP_NAME}" .

echo "==> Copying resources"
cp build/Info.plist "${CONTENTS_DIR}/Info.plist"
if [ -f ui/icon.icns ]; then
    cp ui/icon.icns "${RESOURCES_DIR}/icon.icns"
fi

echo "==> Creating DMG"
# 드래그 설치용 임시 폴더 구성
mkdir -p dist/dmg
cp -r "${BUNDLE_DIR}" dist/dmg/
ln -sf /Applications dist/dmg/Applications

hdiutil create -volname "${APP_NAME}" \
    -srcfolder dist/dmg \
    -ov -format UDZO \
    "dist/${APP_NAME}.dmg"

rm -rf dist/dmg

echo "==> Done: dist/${APP_NAME}.dmg"
