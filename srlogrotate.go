package srlogrotate

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func NewLogger(fileBaseName string) io.Writer {
	return &logger{
		fileBaseName: fileBaseName,
		timeFormat:   "20060102",
		nowFunc:      time.Now,
	}
}

type logger struct {
	fileBaseName string
	timeFormat   string
	nowFunc      func() time.Time

	file *os.File
	mu   sync.Mutex
}

func (l *logger) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file == nil {
		if err := l.openExistingOrNew(); err != nil {
			return 0, err
		}
	}

	if l.shouldRotate() {
		if err := l.rotate(); err != nil {
			return 0, err
		}
	}

	return l.file.Write(p)
}

func (l *logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.close()
}

func (l *logger) close() error {
	if l.file == nil {
		return nil
	}
	err := l.file.Close()
	l.file = nil
	return err
}

func (l *logger) shouldRotate() bool {
	if l.file == nil {
		return false
	}
	if l.file.Name() == l.fileName() {
		return false
	}
	return true
}

func (l *logger) Rotate() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.rotate()
}

func (l *logger) rotate() error {
	if err := l.close(); err != nil {
		return err
	}
	if err := l.openNew(); err != nil {
		return err
	}
	return nil
}

func (l *logger) openExistingOrNew() error {
	fileName := l.fileName()

	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return l.openNew()
	}
	if err != nil {
		return fmt.Errorf("error getting log file info: %s", err)
	}

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return l.openNew()
	}

	l.file = file
	return nil
}

func (l *logger) openNew() error {
	if err := os.MkdirAll(filepath.Dir(l.fileName()), 0755); err != nil {
		return fmt.Errorf("can't make directories for new logfile: %s", err)
	}

	name := l.fileName()
	mode := os.FileMode(0600)
	info, err := os.Stat(name)
	if err == nil {
		mode = info.Mode()
	}

	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("can't open new logfile: %s", err)
	}
	l.file = f

	return nil
}

func (l *logger) fileName() string {
	now := l.nowFunc()
	return fmt.Sprintf("%s.%s", l.fileBaseName, l.timestampStr(now))
}

func (l *logger) timestampStr(tm time.Time) string {
	return tm.Format(l.timeFormat)
}
