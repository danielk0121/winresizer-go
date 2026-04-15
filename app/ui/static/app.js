// ── 언어 설정 ──────────────────────────────────────────────────
const LANG = {
    ko: {
        title:        'WinResizer 설정',
        saveBtn:      '저장 및 적용',
        clearBtn:     '전체 단축키 삭제',
        saveDone:     '저장 완료! 단축키가 즉시 반영되었습니다.',
        saveFail:     '저장 실패.',
        customSection:'커스텀 비율 창 조절',
        customDesc:   '비율(1~100%)을 입력하고 단축키를 설정하세요.',
        dirLeft:      '좌측',
        dirRight:     '우측',
        dirTop:       '상단',
        dirBottom:    '하단',
        hotkeySection:'단축키',
        hotkeyDefault:'단축키 입력',
        hotkeyWaiting:'키 입력 대기...',
        gapSection:   '창 간격 (Gap)',
        confirmClear: '전체 단축키를 삭제할까요?',
        clearDone:    '모든 단축키가 삭제되었습니다. (저장 버튼을 눌러야 반영됩니다)',
        resetBtn:     '기본 세팅값으로 초기화',
        confirmReset: '모든 설정을 기본값으로 초기화할까요?',
        resetDone:    '기본값으로 초기화되었습니다.',
        resetFail:    '초기화 실패.',
        pctError:     '비율은 1~100 사이 정수를 입력하세요.',
        duplicateError: (names) => `단축키 충돌: ${names} 와(과) 동일한 단축키입니다. 저장할 수 없습니다.`,
        applyDone:    (dir, pct) => `${dir} ${pct}% 적용 완료`,
        applyFail:    '적용 실패',
        statusChecking: '권한 확인 중...',
        statusGranted:  '모든 권한 승인됨',
        statusPartial:  '권한 일부 필요 (클릭하여 설정)',
        statusDenied:   '권한 설정 필요 (클릭하여 설정)',
        guideTitle:     'WinResizer 시작하기',
        guideStep1:     '1. <b>손쉬운 사용</b> 권한을 허용해 주세요. <button class="step-btn" onclick="openAccessibilitySettings()">설정 열기</button>',
        guideStep2:     '2. <b>입력 모니터링</b> 권한을 허용해 주세요. <button class="step-btn" onclick="openInputMonitoring()">설정 열기</button>',
        guideStep3:     '3. 설정을 모두 마쳤다면 <b>앱을 종료 후 다시 실행</b>해 주세요.',
        guideStep4:     '4. 모든 준비가 완료되었습니다!',
        guideBtn:       '시스템 설정 열기',
        guideNotice:    '권한 승인 후 이 창이 자동으로 닫히지 않으면 페이지를 새로고침 하세요.',
    },
    en: {
        title:        'WinResizer Settings',
        saveBtn:      'Save & Apply',
        clearBtn:     'Delete All Hotkeys',
        saveDone:     'Saved! Hotkeys applied immediately.',
        saveFail:     'Save failed.',
        customSection:'Custom Ratio Resize',
        customDesc:   'Enter a ratio (1~100%) and set a hotkey.',
        dirLeft:      'Left',
        dirRight:     'Right',
        dirTop:       'Top',
        dirBottom:    'Bottom',
        hotkeySection:'Hotkeys',
        hotkeyDefault:'Press hotkey',
        hotkeyWaiting:'Waiting for key...',
        gapSection:   'Window Gap',
        confirmClear: 'Delete all hotkeys?',
        clearDone:    'All hotkeys deleted. (Must click Save to apply)',
        resetBtn:     'Reset to Defaults',
        confirmReset: 'Reset all settings to default values?',
        resetDone:    'Reset to default settings.',
        resetFail:    'Reset failed.',
        pctError:     'Enter an integer between 1 and 100.',
        duplicateError: (names) => `Hotkey conflict: same as ${names}. Cannot save.`,
        applyDone:    (dir, pct) => `${dir} ${pct}% applied`,
        applyFail:    'Apply failed',
        statusChecking: 'Checking...',
        statusGranted:  'All Permissions Granted',
        statusPartial:  'Permissions Required (Click)',
        statusDenied:   'Permissions Required (Click)',
        guideTitle:     'Getting Started with WinResizer',
        guideStep1:     '1. Enable <b>Accessibility</b> permission. <button class="step-btn" onclick="openAccessibilitySettings()">Open Settings</button>',
        guideStep2:     '2. Enable <b>Input Monitoring</b> permission. <button class="step-btn" onclick="openInputMonitoring()">Open Settings</button>',
        guideStep3:     '3. Once done, <b>restart the app</b> for changes to take effect.',
        guideStep4:     '4. Everything is ready!',
        guideBtn:       'Open System Settings',
        guideNotice:    'If this window doesn\'t close after granting, please refresh the page.',
    },
};

let currentLang = localStorage.getItem('lang') || 'ko';

function t(key, ...args) {
    const val = LANG[currentLang][key];
    return typeof val === 'function' ? val(...args) : val;
}

function setLang(lang) {
    currentLang = lang;
    localStorage.setItem('lang', lang);
    applyLang();
    // 가이드 내부 HTML은 <b> 태그 등이 포함될 수 있어 innerHTML로 처리
    updateGuideUI();
    // 단축키 버튼 텍스트 갱신 (녹화 중이 아닌 버튼만)
    renderHotkeys();
    renderCustomHotkeys();
}

function applyLang() {
    // data-i18n 속성 요소 텍스트 일괄 교체
    document.querySelectorAll('[data-i18n]').forEach(el => {
        const key = el.getAttribute('data-i18n');
        if (LANG[currentLang][key] !== undefined) {
            // 가이드 단계나 강조 문구 등 HTML이 포함된 경우 innerHTML 사용
            if (key.startsWith('guideStep')) {
                el.innerHTML = t(key);
            } else {
                el.textContent = t(key);
            }
        }
    });
    // 언어 버튼 활성화 표시
    document.getElementById('lang-ko').classList.toggle('active', currentLang === 'ko');
    document.getElementById('lang-en').classList.toggle('active', currentLang === 'en');
    // html lang 속성 갱신
    document.documentElement.lang = currentLang;
}

function updateGuideUI() {
    // 가이드 단계를 다시 렌더링하여 언어 변경 반영
    const steps = ['guideStep1', 'guideStep2', 'guideStep3', 'guideStep4'];
    steps.forEach((key, index) => {
        const el = document.getElementById(`step${index + 1}`);
        if (el) {
            const textEl = el.querySelector('.step-text');
            if (textEl) textEl.innerHTML = t(key);
        }
    });
}

// ── 커스텀 비율 키 설정 ────────────────────────────────────────
const CUSTOM_KEYS = ['Left Custom', 'Right Custom', 'Top Custom', 'Bottom Custom'];

const CUSTOM_PCT_IDS = {
    'Left Custom': 'pct-left',
    'Right Custom': 'pct-right',
    'Top Custom': 'pct-top',
    'Bottom Custom': 'pct-bottom',
};

// 단축키 항목 한글 이름 맵
const HOTKEY_LABELS = {
    ko: {
        'Left':             '좌측 1/2',
        'Right':            '우측 1/2',
        'Top':              '상단 1/2',
        'Bottom':           '하단 1/2',
        'Left 1/3':         '좌측 1/3',
        'Center 1/3':       '중앙 1/3',
        'Right 1/3':        '우측 1/3',
        'Left 2/3':         '좌측 2/3',
        'Right 2/3':        '우측 2/3',
        'Top Left 1/4':     '좌상단 1/4',
        'Top Right 1/4':    '우상단 1/4',
        'Bottom Left 1/4':  '좌하단 1/4',
        'Bottom Right 1/4': '우하단 1/4',
        'Maximize':         '최대화',
        'Restore':          '복구',
    },
    en: {},  // 영어는 키 이름 그대로 사용
};

function hotkeyLabel(name) {
    return HOTKEY_LABELS[currentLang]?.[name] || name;
}

// 단축키 섹션 렌더링 순서
const HOTKEY_ORDER = [
    'Left', 'Right', 'Top', 'Bottom',
    'Left 1/3', 'Center 1/3', 'Right 1/3',
    'Left 2/3', 'Right 2/3',
    'Top Left 1/4', 'Top Right 1/4', 'Bottom Left 1/4', 'Bottom Right 1/4',
    'Maximize', 'Restore',
];

let config = {};
let initialConfig = {};
let recordingKey = null;

function deepCopy(obj) {
    return JSON.parse(JSON.stringify(obj));
}

async function loadConfig() {
    const res = await fetch('/api/config');
    config = await res.json();
    initialConfig = deepCopy(config);
    loadConfigUI();
}

function loadConfigUI() {
    // 1. 모든 DOM 값(input)들을 먼저 설정
    document.getElementById('gap').value = config.settings?.gap ?? 5;
    
    for (const name of CUSTOM_KEYS) {
        const info = config.shortcuts?.[name];
        if (!info) continue;
        const pctId = CUSTOM_PCT_IDS[name];
        const mode = info.mode || '';
        const match = mode.match(/_custom:(\d+)$/);
        if (match && pctId) {
            document.getElementById(pctId).value = match[1];
        }
        const row = document.getElementById('row-' + name);
        if (row) row.dataset.initialized = "true";
    }

    // 2. 그 다음 UI 렌더링 및 음영 판단 수행
    applyLang();
    renderHotkeys();
    renderCustomHotkeys();
    renderGap();
}

function isModified(name) {
    const cur = config.shortcuts?.[name];
    const ini = initialConfig.shortcuts?.[name];
    if (!cur || !ini) return false;

    // keycode/modifiers 비교
    if (cur.keycode !== ini.keycode || cur.modifiers !== ini.modifiers) return true;

    // 커스텀 비율 모드인 경우 비율 비교
    if (CUSTOM_KEYS.includes(name)) {
        const pctId = CUSTOM_PCT_IDS[name];
        const inputEl = document.getElementById(pctId);
        if (!inputEl) return false;
        const curPct = inputEl.value;
        const mode = ini.mode || '';
        const match = mode.match(/_custom:(\d+)$/);
        const iniPct = match ? match[1] : '';
        if (String(curPct) !== String(iniPct)) return true;
    }

    return false;
}

function isGapModified() {
    const curGap = document.getElementById('gap').value;
    const iniGap = initialConfig.settings?.gap ?? 5;
    return parseInt(curGap) !== iniGap;
}

function renderHotkeys() {
    const list = document.getElementById('hotkey-list');
    list.innerHTML = '';
    const shortcuts = config.shortcuts || {};
    for (const name of HOTKEY_ORDER) {
        const info = shortcuts[name];
        if (!info) continue;
        const display = info.keycode ? info.display : t('hotkeyDefault');
        const row = document.createElement('div');
        row.className = 'row' + (isModified(name) ? ' modified' : '');
        row.id = 'row-' + name;
        row.innerHTML = `
            <span class="label">${hotkeyLabel(name)}</span>
            <button class="hotkey-btn" id="btn-${name}" onclick="startRecording('${name}')">${display}</button>
            <button class="delete-btn" onclick="deleteHotkey('${name}')">✕</button>
        `;
        list.appendChild(row);
    }
}

function renderCustomHotkeys() {
    for (const name of CUSTOM_KEYS) {
        const info = config.shortcuts?.[name];
        if (!info) continue;
        
        const row = document.getElementById('row-' + name);
        if (row) {
            row.className = 'custom-row' + (isModified(name) ? ' modified' : '');
        }

        const btn = document.getElementById('btn-' + name);
        if (btn && !btn.classList.contains('recording')) {
            btn.textContent = info.keycode ? info.display : t('hotkeyDefault');
        }
    }
}

function renderGap() {
    const row = document.getElementById('row-gap');
    if (row) {
        row.className = 'gap-row' + (isGapModified() ? ' modified' : '');
    }
}

function startRecording(keyName) {
    if (recordingKey) stopRecording();
    recordingKey = keyName;
    const btn = document.getElementById('btn-' + keyName);
    btn.textContent = t('hotkeyWaiting');
    btn.classList.add('recording');
}

function stopRecording() {
    if (!recordingKey) return;
    const btn = document.getElementById('btn-' + recordingKey);
    if (btn) btn.classList.remove('recording');
    recordingKey = null;
    
    // 변경 사항이 있을 수 있으므로 렌더링 갱신
    renderHotkeys();
    renderCustomHotkeys();
}

// ── Carbon keycode 변환 맵 (브라우저 e.code → macOS Virtual Key Code) ──
// https://developer.apple.com/documentation/carbon/1805242-virtual_key_codes
const KEYCODE_MAP = {
    'KeyA':0,'KeyS':1,'KeyD':2,'KeyF':3,'KeyH':4,'KeyG':5,'KeyZ':6,'KeyX':7,
    'KeyC':8,'KeyV':9,'KeyB':11,'KeyQ':12,'KeyW':13,'KeyE':14,'KeyR':15,
    'KeyY':16,'KeyT':17,'Digit1':18,'Digit2':19,'Digit3':20,'Digit4':21,
    'Digit6':22,'Digit5':23,'Equal':24,'Digit9':25,'Digit7':26,'Minus':27,
    'Digit8':28,'Digit0':29,'BracketRight':30,'KeyO':31,'KeyU':32,
    'BracketLeft':33,'KeyI':34,'KeyP':35,'Enter':36,'KeyL':37,'KeyJ':38,
    'Quote':39,'KeyK':40,'Semicolon':41,'Backslash':42,'Comma':43,'Slash':44,
    'KeyN':45,'KeyM':46,'Period':47,'Tab':48,'Space':49,'Backquote':50,
    'Backspace':51,'Escape':53,'ArrowLeft':123,'ArrowRight':124,
    'ArrowDown':125,'ArrowUp':126,'F1':122,'F2':120,'F3':99,'F4':118,
    'F5':96,'F6':97,'F7':98,'F8':100,'F9':101,'F10':109,'F11':103,'F12':111,
};

// Carbon modifier 비트 플래그
const MOD_SHIFT   = 1 << 9;  // shiftKey
const MOD_CTRL    = 1 << 12; // ctrlKey
const MOD_OPT     = 1 << 11; // altKey (Option)
const MOD_CMD     = 1 << 8;  // metaKey (Command)

document.addEventListener('keydown', (e) => {
    if (!recordingKey) return;
    e.preventDefault();

    // Backspace/Delete: 단축키 삭제
    if (e.key === 'Backspace' || e.key === 'Delete') {
        config.shortcuts[recordingKey].keycode   = 0;
        config.shortcuts[recordingKey].modifiers = 0;
        config.shortcuts[recordingKey].display   = '';
        stopRecording();
        return;
    }

    // 수식키만 눌린 경우 무시
    if (['Control', 'Alt', 'Meta', 'Shift'].includes(e.key)) return;

    const keycode = KEYCODE_MAP[e.code];
    if (keycode === undefined) return; // 미지원 키 무시

    // 수식키 조합 (최소 1개 이상 필요)
    if (!e.ctrlKey && !e.altKey && !e.metaKey && !e.shiftKey) return;

    let modifiers = 0;
    if (e.ctrlKey)  modifiers |= MOD_CTRL;
    if (e.altKey)   modifiers |= MOD_OPT;
    if (e.metaKey)  modifiers |= MOD_CMD;
    if (e.shiftKey) modifiers |= MOD_SHIFT;

    // 사람이 읽기 좋은 표시용 문자열
    const parts = [];
    if (e.ctrlKey)  parts.push('ctrl');
    if (e.altKey)   parts.push('opt');
    if (e.metaKey)  parts.push('cmd');
    if (e.shiftKey) parts.push('shift');
    parts.push(e.key.length === 1 ? e.key.toUpperCase() : e.key);
    const display = parts.join(' + ');

    config.shortcuts[recordingKey].keycode   = keycode;
    config.shortcuts[recordingKey].modifiers = modifiers;
    config.shortcuts[recordingKey].display   = display;
    stopRecording();
});

function getTimeString() {
    const now = new Date();
    const tzo = -now.getTimezoneOffset();
    const dif = tzo >= 0 ? '+' : '-';
    const pad = (num) => String(num).padStart(2, '0');
    
    const iso = now.getFullYear() +
        '-' + pad(now.getMonth() + 1) +
        '-' + pad(now.getDate()) +
        'T' + pad(now.getHours()) +
        ':' + pad(now.getMinutes()) +
        ':' + pad(now.getSeconds()) +
        '.' + String(now.getMilliseconds()).padStart(3, '0') +
        dif + pad(Math.floor(Math.abs(tzo) / 60)) +
        ':' + pad(Math.abs(tzo) % 60);
    return `[${iso}]\n`;
}

function deleteHotkey(name) {
    if (!config.shortcuts[name]) return;
    config.shortcuts[name].keycode   = 0;
    config.shortcuts[name].modifiers = 0;
    config.shortcuts[name].display   = '';
    renderHotkeys();
    renderCustomHotkeys();
}

function clearAll() {
    if (!confirm(t('confirmClear'))) return;
    for (const name of Object.keys(config.shortcuts || {})) {
        config.shortcuts[name].keycode   = 0;
        config.shortcuts[name].modifiers = 0;
        config.shortcuts[name].display   = '';
    }
    renderHotkeys();
    renderCustomHotkeys();
    renderGap();

    const status = document.getElementById('status');
    status.style.color = '#f39c12';
    status.textContent = getTimeString() + t('clearDone');
}

async function saveConfig() {
    const status = document.getElementById('status');

    // 비율 유효성 검사 (1~100 정수)
    const dirMap = {
        'Left Custom': 'left',
        'Right Custom': 'right',
        'Top Custom': 'top',
        'Bottom Custom': 'bottom',
    };
    for (const [name, dir] of Object.entries(dirMap)) {
        const pctId = CUSTOM_PCT_IDS[name];
        const raw = document.getElementById(pctId).value;
        const pct = parseInt(raw);
        if (raw === '' || isNaN(pct) || pct < 1 || pct > 100 || String(pct) !== raw.trim()) {
            status.style.color = '#e74c3c';
            status.textContent = getTimeString() + t('pctError');
            return;
        }
    }

    // 중복 단축키 검사 (keycode + modifiers 조합 기준)
    const keyCombToNames = {};
    for (const [name, info] of Object.entries(config.shortcuts || {})) {
        if (!info.keycode) continue;
        const key = `${info.keycode}:${info.modifiers}`;
        if (!keyCombToNames[key]) keyCombToNames[key] = [];
        keyCombToNames[key].push(name);
    }
    const conflicts = Object.values(keyCombToNames).filter(names => names.length > 1);
    if (conflicts.length > 0) {
        const conflictDesc = conflicts.map(names => names.map(hotkeyLabel).join(' ↔ ')).join(', ');
        status.style.color = '#e74c3c';
        status.textContent = getTimeString() + t('duplicateError', conflictDesc);
        return;
    }

    config.settings = config.settings || {};
    config.settings.gap = parseInt(document.getElementById('gap').value) || 0;

    for (const [name, dir] of Object.entries(dirMap)) {
        const pctId = CUSTOM_PCT_IDS[name];
        const pct = parseInt(document.getElementById(pctId).value);
        if (config.shortcuts[name]) {
            config.shortcuts[name].mode = `${dir}_custom:${pct}`;
        }
    }

    const res = await fetch('/api/config', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(config)
    });
    if (res.ok) {
        initialConfig = deepCopy(config); // 저장 성공 시 초기값 갱신하여 음영 제거
        renderHotkeys();
        renderCustomHotkeys();
        renderGap();
        status.style.color = '#2ecc71';
        status.textContent = getTimeString() + t('saveDone');
    } else {
        status.style.color = '#e74c3c';
        status.textContent = getTimeString() + t('saveFail');
    }
}

async function resetConfig() {
    if (!confirm(t('confirmReset'))) return;
    const res = await fetch('/api/config/reset', { method: 'POST' });
    const status = document.getElementById('status');
    if (res.ok) {
        const data = await res.json();
        config = data.config; // 로컬 설정만 업데이트 (initialConfig는 그대로 두어 음영 발생 유도)
        loadConfigUI();
        status.style.color = '#f39c12';
        status.textContent = getTimeString() + t('resetDone') + ' (저장 버튼을 눌러야 반영됩니다)';
    } else {
        status.style.color = '#e74c3c';
        status.textContent = getTimeString() + t('resetFail');
    }
}

loadConfig();

// ── 상태 체크 ──────────────────────────────────────────────────
async function checkStatus() {
    const badge = document.getElementById('status-badge');
    if (!badge) return;

    try {
        const res = await fetch('/api/status');
        const data = await res.json();
        const acc = data.accessibility_granted;
        const inp = data.input_monitoring_granted;

        if (acc && inp) {
            badge.className = 'status-badge status-granted';
            badge.innerHTML = `<span>${t('statusGranted')}</span>`;
            badge.onclick = null;
            // 모든 권한 승인 — 가이드 오버레이 숨김
            hideGuideOverlay();
        } else {
            badge.className = 'status-badge status-denied';
            const statusKey = (acc || inp) ? 'statusPartial' : 'statusDenied';
            badge.innerHTML = `<span>${t(statusKey)}</span>`;
            // 권한 미승인 — 가이드 오버레이 표시
            showGuideOverlay(acc, inp);
        }
    } catch (e) {
        console.error('Status check failed:', e);
    }
}

function showGuideOverlay(acc, inp) {
    const overlay = document.getElementById('guide-overlay');
    if (!overlay) return;
    overlay.style.display = 'flex';

    // 완료된 단계 체크 표시
    const step1 = document.getElementById('step1');
    const step2 = document.getElementById('step2');
    const step3 = document.getElementById('step3');
    const step4 = document.getElementById('step4');
    const mainBtn = document.getElementById('guide-main-btn');

    if (step1) step1.className = 'guide-step' + (acc ? ' step-done' : '');
    if (step2) step2.className = 'guide-step' + (inp ? ' step-done' : '');

    // 두 권한 모두 승인 시 재실행 안내 단계 표시
    if (acc && inp) {
        if (step3) step3.style.display = '';
        if (step4) step4.style.display = '';
        if (mainBtn) mainBtn.style.display = 'none';
    } else {
        if (step3) step3.style.display = '';
        if (step4) step4.style.display = 'none';
        // 미승인 권한에 맞게 메인 버튼 동작 변경
        if (mainBtn) {
            mainBtn.style.display = '';
            if (!acc) {
                mainBtn.onclick = openAccessibilitySettings;
            } else {
                mainBtn.onclick = openInputMonitoring;
            }
        }
    }
}

function hideGuideOverlay() {
    const overlay = document.getElementById('guide-overlay');
    if (overlay) overlay.style.display = 'none';
}

async function openAccessibilitySettings() {
    await fetch('/api/execute', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ mode: 'open_accessibility' })
    });
}

async function openInputMonitoring() {
    await fetch('/api/execute', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ mode: 'open_input_monitoring' })
    });
}

// 5초마다 권한 상태 체크
setInterval(checkStatus, 5000);
setTimeout(checkStatus, 500);
