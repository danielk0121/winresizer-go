#!/bin/bash
# SVG → PNG 각 해상도 생성 → iconutil로 icns 변환
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SVG="${SCRIPT_DIR}/wr_icon.svg"
ICONSET="${SCRIPT_DIR}/icon.iconset"
APP_UI="${SCRIPT_DIR}/../../app/ui"

mkdir -p "${ICONSET}"

# iconset 규격 해상도
sizes=(16 32 64 128 256 512 1024)
declare -A filenames=(
    [16]="icon_16x16.png"
    [32]="icon_16x16@2x.png icon_32x32.png"
    [64]="icon_32x32@2x.png"
    [128]="icon_128x128.png"
    [256]="icon_128x128@2x.png icon_256x256.png"
    [512]="icon_256x256@2x.png icon_512x512.png"
    [1024]="icon_512x512@2x.png"
)

echo "==> SVG → PNG 변환"
for size in "${sizes[@]}"; do
    tmp="${SCRIPT_DIR}/tmp_${size}.png"
    rsvg-convert -w "${size}" -h "${size}" "${SVG}" -o "${tmp}"
    for fname in ${filenames[$size]}; do
        cp "${tmp}" "${ICONSET}/${fname}"
    done
    rm "${tmp}"
done

echo "==> iconutil → icon.icns"
iconutil -c icns "${ICONSET}" -o "${APP_UI}/icon.icns"

echo "==> tray_icon.png (22px)"
rsvg-convert -w 22 -h 22 "${SVG}" -o "${APP_UI}/tray_icon.png"

echo "==> favicon.png"
cp "${APP_UI}/tray_icon.png" "${APP_UI}/static/favicon.png"

echo "==> 완료"
echo "  icon.icns: ${APP_UI}/icon.icns"
echo "  tray_icon.png: ${APP_UI}/tray_icon.png"
echo "  favicon.png: ${APP_UI}/static/favicon.png"
