function setLanguage(lang) {
    document.documentElement.setAttribute('lang', lang);
    localStorage.setItem('preferred-lang', lang);

    const btnKo = document.getElementById('btn-ko');
    const btnEn = document.getElementById('btn-en');
    if (lang === 'ko') {
        btnKo.classList.add('active');
        btnEn.classList.remove('active');
        document.title = '윈리사이저 - macOS를 위한 가볍고 강력한 창 관리자';
    } else {
        btnEn.classList.add('active');
        btnKo.classList.remove('active');
        document.title = 'WinResizer - Lightweight & Powerful Window Manager for macOS';
    }
}

window.addEventListener('DOMContentLoaded', () => {
    const cachedLang = localStorage.getItem('preferred-lang');
    if (cachedLang === 'ko' || cachedLang === 'en') {
        setLanguage(cachedLang);
        return;
    }

    const userLang = navigator.language || navigator.userLanguage;
    setLanguage(userLang.startsWith('ko') ? 'ko' : 'en');
});
