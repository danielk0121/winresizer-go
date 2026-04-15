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

echo "==> Building Universal Binary"
GOARCH=amd64 go build -o dist/winresizer-amd64 .
GOARCH=arm64 go build -o dist/winresizer-arm64 .
lipo -create dist/winresizer-amd64 dist/winresizer-arm64 -output "${MACOS_DIR}/${APP_NAME}"
rm dist/winresizer-amd64 dist/winresizer-arm64

echo "==> Copying resources"
cp build/Info.plist "${CONTENTS_DIR}/Info.plist"
if [ -f ui/icon.icns ]; then
    cp ui/icon.icns "${RESOURCES_DIR}/icon.icns"
fi

echo "==> Creating DMG"
hdiutil create -volname "${APP_NAME}" \
    -srcfolder dist/ \
    -ov -format UDZO \
    "dist/${APP_NAME}.dmg"

echo "==> Done: dist/${APP_NAME}.dmg"
