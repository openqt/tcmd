package operations

import (
	"fmt"
	"sync"
	"time"
)

// LogEntry records one file operation.
type LogEntry struct {
	Time    time.Time
	Message string
}

// Log stores operation history.
type Log struct {
	mu      sync.Mutex
	entries []LogEntry
}

// Global operation log.
var DefaultLog Log

// Add appends a log line.
func (l *Log) Add(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = append(l.entries, LogEntry{Time: time.Now(), Message: msg})
}

// Entries returns a copy of log lines.
func (l *Log) Entries() []LogEntry {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]LogEntry, len(l.entries))
	copy(out, l.entries)
	return out
}

// Format returns log as text.
func (l *Log) Format() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	var s string
	for _, e := range l.entries {
		s += fmt.Sprintf("%s %s\n", e.Time.Format("15:04:05"), e.Message)
	}
	return s
}
