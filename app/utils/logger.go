package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// 전역 로거
var Log *Logger

// Logger는 파일 + 콘솔 동시 출력을 지원하는 로거입니다.
type Logger struct {
	logger  *log.Logger
	logFile *os.File
}

func init() {
	var err error
	Log, err = newLogger()
	if err != nil {
		// 로거 초기화 실패 시 콘솔 전용으로 폴백
		Log = &Logger{logger: log.New(os.Stdout, "", 0)}
		Log.Errorf("로거 초기화 실패: %v", err)
	}
}

func newLogger() (*Logger, error) {
	logDir := filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "WinResizer", "log")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("로그 디렉토리 생성 실패: %w", err)
	}

	// 날짜별 롤링: winresizer_2026-04-15.log
	today := time.Now().Format("2006-01-02")
	logPath := filepath.Join(logDir, fmt.Sprintf("winresizer_%s.log", today))

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("로그 파일 열기 실패: %w", err)
	}

	// 파일 + 콘솔 동시 출력
	multi := io.MultiWriter(f, os.Stdout)
	l := log.New(multi, "", 0)

	return &Logger{logger: l, logFile: f}, nil
}

func (l *Logger) format(level, msg string) string {
	now := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("%s [%s] %s", now, level, msg)
}

func (l *Logger) Debugf(format string, args ...any) {
	l.logger.Println(l.format("DEBUG", fmt.Sprintf(format, args...)))
}

func (l *Logger) Infof(format string, args ...any) {
	l.logger.Println(l.format("INFO", fmt.Sprintf(format, args...)))
}

func (l *Logger) Warnf(format string, args ...any) {
	l.logger.Println(l.format("WARN", fmt.Sprintf(format, args...)))
}

func (l *Logger) Errorf(format string, args ...any) {
	l.logger.Println(l.format("ERROR", fmt.Sprintf(format, args...)))
}

// Close는 로그 파일을 닫습니다. 앱 종료 시 호출하세요.
func (l *Logger) Close() {
	if l.logFile != nil {
		l.logFile.Close()
	}
}
