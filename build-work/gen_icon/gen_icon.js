// gen_icon.js: SVG → PNG 각 해상도 생성 → iconutil로 icns 변환
const sharp = require('sharp');
const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

const SCRIPT_DIR = __dirname;
const SVG_PATH = path.join(SCRIPT_DIR, 'wr_icon.svg');
const ICONSET_DIR = path.join(SCRIPT_DIR, 'out', 'icon.iconset');
const APP_UI = path.join(SCRIPT_DIR, '../../app/ui');

// iconset 규격
const specs = [
  { size: 16,   filename: 'icon_16x16.png' },
  { size: 32,   filename: 'icon_16x16@2x.png' },
  { size: 32,   filename: 'icon_32x32.png' },
  { size: 64,   filename: 'icon_32x32@2x.png' },
  { size: 128,  filename: 'icon_128x128.png' },
  { size: 256,  filename: 'icon_128x128@2x.png' },
  { size: 256,  filename: 'icon_256x256.png' },
  { size: 512,  filename: 'icon_256x256@2x.png' },
  { size: 512,  filename: 'icon_512x512.png' },
  { size: 1024, filename: 'icon_512x512@2x.png' },
];

async function main() {
  fs.mkdirSync(ICONSET_DIR, { recursive: true });

  const svgBuf = fs.readFileSync(SVG_PATH);

  console.log('==> SVG → PNG 변환');
  for (const spec of specs) {
    const outPath = path.join(ICONSET_DIR, spec.filename);
    await sharp(svgBuf)
      .resize(spec.size, spec.size)
      .png()
      .toFile(outPath);
    console.log(`  ${spec.filename} (${spec.size}px)`);
  }

  console.log('==> iconutil → icon.icns');
  const icnsOut = path.join(APP_UI, 'icon.icns');
  execSync(`iconutil -c icns "${ICONSET_DIR}" -o "${icnsOut}"`);
  console.log(`  icon.icns: ${icnsOut}`);

  console.log('==> tray_icon.png (22px)');
  const trayLocal = path.join(SCRIPT_DIR, 'out', 'tray_icon.png');
  const trayOut = path.join(APP_UI, 'tray_icon.png');
  // 작은 사이즈는 먼저 큰 해상도로 렌더링 후 다운스케일
  const trayBuf = await sharp(svgBuf).resize(88, 88).png().toBuffer();
  await sharp(trayBuf).resize(22, 22).png().toFile(trayLocal);
  fs.copyFileSync(trayLocal, trayOut);
  console.log(`  tray_icon.png: ${trayOut}`);

  console.log('==> favicon.png');
  const faviconOut = path.join(APP_UI, 'static', 'favicon.png');
  fs.copyFileSync(trayOut, faviconOut);
  console.log(`  favicon.png: ${faviconOut}`);

  console.log('==> 완료');
}

main().catch(err => { console.error(err); process.exit(1); });
