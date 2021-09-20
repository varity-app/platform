package logging

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/spf13/viper"
)

// For more details regarding log levels, see link below:
// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogSeverity
const (
	// SeverityCritical means events cause more severe problems or outages.
	SeverityCritical = "CRITICAL"
	// SeverityError means events are likely to cause problems.
	SeverityError = "ERROR"
	// SeverityWarning means events might cause problems.
	SeverityWarning = "WARNING"
	// SeverityNotice means normal but significant events,
	// such as start up, shut down, or a configuration change.
	SeverityNotice = "NOTICE"
	// SeverityInfo means routine information, such as ongoing status or performance.
	SeverityInfo = "INFO"
	// SeverityDebug means debug or trace information.
	SeverityDebug = "DEBUG"
)

// Map a log level to a numerical representation. This makes it easy to compare levels.
var levelMap = map[string]int{
	SeverityCritical: 600,
	SeverityError:    500,
	SeverityWarning:  400,
	SeverityNotice:   300,
	SeverityInfo:     200,
	SeverityDebug:    100,
}

// LogEntry is a structured log entry.
type LogEntry struct {
	Timestamp   time.Time `json:"timestamp"`
	Severity    string    `json:"severity"`
	TextPayload string    `json:"textPayload"`
	Source      string    `json:"source"`
}

// Logger is a structured, leveled logger.
type Logger struct {
	minLevel int
}

// NewLogger constructs a new structured Logger
func NewLogger() *Logger {
	minLevel := viper.GetString("logging.level")
	return &Logger{minLevel: levelMap[minLevel]}
}

// Fatal is a writes a log with level CRITICAL and also exits the program.
func (l *Logger) Fatal(err error) {
	// Create log entry
	entry := LogEntry{
		Severity:    SeverityCritical,
		TextPayload: err.Error(),
		Timestamp:   time.Now(),
		Source:      getSource(),
	}

	// Serialize entry
	serializedEntry, err := json.Marshal(entry)
	if err != nil {
		log.Fatalf("logging.SerializeLogEntry: %v", err)
	}

	// Write entry to stdout
	fmt.Println(string(serializedEntry))

	// Exit program
	os.Exit(1)
}

// Error writes a log with an ERROR level.
func (l *Logger) Error(err error) {

	// Do nothing if min level is higher than this level.
	if l.minLevel > levelMap[SeverityError] {
		return
	}

	// Create log entry
	entry := LogEntry{
		Severity:    SeverityError,
		TextPayload: err.Error(),
		Timestamp:   time.Now(),
		Source:      getSource(),
	}

	// Serialize entry
	serializedEntry, err := json.Marshal(entry)
	if err != nil {
		log.Fatalf("logging.SerializeLogEntry: %v", err)
	}

	// Write entry to stdout
	fmt.Println(string(serializedEntry))
}

// Warning writes a log with a WARNING level.
func (l *Logger) Warning(message string) {

	// Do nothing if min level is higher than this level.
	if l.minLevel > levelMap[SeverityWarning] {
		return
	}

	// Create log entry
	entry := LogEntry{
		Severity:    SeverityWarning,
		TextPayload: message,
		Timestamp:   time.Now(),
		Source:      getSource(),
	}

	// Serialize entry
	serializedEntry, err := json.Marshal(entry)
	if err != nil {
		log.Fatalf("logging.SerializeLogEntry: %v", err)
	}

	// Write entry to stdout
	fmt.Println(string(serializedEntry))
}

// Info writes a log with an INFO level.
func (l *Logger) Info(message string) {

	// Do nothing if min level is higher than this level.
	if l.minLevel > levelMap[SeverityInfo] {
		return
	}

	// Create log entry
	entry := LogEntry{
		Severity:    SeverityInfo,
		TextPayload: message,
		Timestamp:   time.Now(),
		Source:      getSource(),
	}

	// Serialize entry
	serializedEntry, err := json.Marshal(entry)
	if err != nil {
		log.Fatalf("logging.SerializeLog: %v", err)
	}

	// Write entry to stdout
	fmt.Println(string(serializedEntry))
}

// Debug writes a log with an DEBUG level.
func (l *Logger) Debug(message string) {

	// Do nothing if min level is higher than this level.
	if l.minLevel > levelMap[SeverityDebug] {
		return
	}

	// Create log entry
	entry := LogEntry{
		Severity:    SeverityDebug,
		TextPayload: message,
		Timestamp:   time.Now(),
		Source:      getSource(),
	}

	// Serialize entry
	serializedEntry, err := json.Marshal(entry)
	if err != nil {
		log.Fatalf("logging.SerializeLog: %v", err)
	}

	// Write entry to stdout
	fmt.Println(string(serializedEntry))
}

// Helper method for getting source file calling the log
func getSource() string {
	_, filename, _, _ := runtime.Caller(2)
	return filename
}
