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
    
    // pynput 비교 (단축키 삭제 포함)
    if (cur.pynput !== ini.pynput) return true;
    
    // 커스텀 비율 모드인 경우 비율 비교
    if (CUSTOM_KEYS.includes(name)) {
        const pctId = CUSTOM_PCT_IDS[name];
        const inputEl = document.getElementById(pctId);
        if (!inputEl) return false;

        const curPct = inputEl.value;
        const mode = ini.mode || '';
        const match = mode.match(/_custom:(\d+)$/);
        const iniPct = match ? match[1] : '';
        
        // 문자열로 비교하여 타입 차이로 인한 이슈 방지
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
        const display = info.pynput ? info.display : t('hotkeyDefault');
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
            btn.textContent = info.pynput ? info.display : t('hotkeyDefault');
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

document.addEventListener('keydown', (e) => {
    if (!recordingKey) return;
    e.preventDefault();

    if (e.key === 'Backspace' || e.key === 'Delete') {
        config.shortcuts[recordingKey].pynput = '';
        config.shortcuts[recordingKey].display = '';
        stopRecording();
        return;
    }

    if (['Control', 'Alt', 'Meta', 'Shift'].includes(e.key)) return;

    // 브라우저 e.key → pynput 키 이름 변환 맵
    const KEY_MAP = {
        'arrowleft':  'left',
        'arrowright': 'right',
        'arrowup':    'up',
        'arrowdown':  'down',
        'enter':      'enter',
        'escape':     'esc',
        'tab':        'tab',
        'space':      'space',
        ' ':          'space',
    };

    const parts = [];
    if (e.ctrlKey)  parts.push('<ctrl>');
    if (e.altKey)   parts.push('<alt>');
    if (e.metaKey)  parts.push('<cmd>');
    if (e.shiftKey) parts.push('<shift>');

    const rawKey = e.key.toLowerCase();
    const mappedKey = KEY_MAP[rawKey] || rawKey;
    const key = mappedKey.length === 1 ? mappedKey : `<${mappedKey}>`;
    parts.push(key);

    const pynput = parts.join('+');
    const display = pynput.replace(/[<>]/g, '').replace(/\+/g, ' + ');

    config.shortcuts[recordingKey].pynput = pynput;
    config.shortcuts[recordingKey].display = display;
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

function clearAll() {
    if (!confirm(t('confirmClear'))) return;
    for (const name of Object.keys(config.shortcuts || {})) {
        config.shortcuts[name].pynput = '';
        config.shortcuts[name].display = '';
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

    // 중복 단축키 검사
    const pynputToNames = {};
    for (const [name, info] of Object.entries(config.shortcuts || {})) {
        const pk = info.pynput;
        if (!pk) continue;
        if (!pynputToNames[pk]) pynputToNames[pk] = [];
        pynputToNames[pk].push(name);
    }
    const conflicts = Object.values(pynputToNames).filter(names => names.length > 1);
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
    const guide = document.getElementById('guide-overlay');
    const mainBtn = document.getElementById('guide-main-btn');
    if (!badge) return;

    try {
        const res = await fetch('/api/status');
        const data = await res.json();
        
        const acc = data.accessibility_granted;
        const inp = data.input_monitoring_granted;

        // 1단계: 손쉬운 사용
        const step1 = document.getElementById('step1');
        if (step1) step1.classList.toggle('completed', acc);

        // 2단계: 입력 모니터링
        const step2 = document.getElementById('step2');
        if (step2) step2.classList.toggle('completed', inp);

        // 3단계: 앱 재실행 (두 권한 모두 있고 앱이 떠있으면 안내)
        const step3 = document.getElementById('step3');
        if (step3) {
            if (acc && inp) {
                step3.style.fontWeight = 'bold';
                step3.style.color = '#2ecc71';
            } else {
                step3.style.fontWeight = 'normal';
                step3.style.color = '#ccc';
            }
        }

        // 전체 배지 업데이트
        if (acc && inp) {
            badge.className = 'status-badge status-granted';
            badge.innerHTML = `<span data-i18n="statusGranted">${t('statusGranted')}</span>`;
            badge.onclick = null;
            if (guide) guide.style.display = 'none';
        } else {
            badge.className = 'status-badge status-denied';
            const statusKey = (acc || inp) ? 'statusPartial' : 'statusDenied';
            badge.innerHTML = `<span data-i18n="${statusKey}">${t(statusKey)}</span>`;
            
            badge.onclick = !acc ? openAccessibilitySettings : openInputMonitoring;
            
            if (guide) {
                guide.style.display = 'flex';
                if (mainBtn) {
                    mainBtn.onclick = !acc ? openAccessibilitySettings : openInputMonitoring;
                }
            }
        }
    } catch (e) {
        console.error('Status check failed:', e);
    }
}

async function openAccessibilitySettings() {
    console.log('Requesting to open accessibility settings...');
    try {
        const res = await fetch('/api/execute', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ mode: 'open_accessibility' })
        });
        if (res.ok) {
            console.log('Accessibility settings request successful');
        } else {
            console.error('Failed to open accessibility settings:', await res.text());
        }
    } catch (e) {
        console.error('Error calling /api/execute:', e);
    }
}

async function openInputMonitoring() {
    console.log('Requesting to open input monitoring settings...');
    try {
        const res = await fetch('/api/execute', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ mode: 'open_input_monitoring' })
        });
        if (res.ok) {
            console.log('Input monitoring settings request successful');
        } else {
            console.error('Failed to open input monitoring settings:', await res.text());
        }
    } catch (e) {
        console.error('Error calling /api/execute:', e);
    }
}

// 5초마다 권한 상태 체크
setInterval(checkStatus, 5000);
setTimeout(checkStatus, 500); // 로드 직후 한 번 실행
