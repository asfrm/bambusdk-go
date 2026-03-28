package logger

import (
	"log"
	"os"
)

var (
	// Debug logger for detailed information
	Debug *log.Logger
	// Info logger for general information
	Info *log.Logger
	// Warn logger for warnings
	Warn *log.Logger
	// Error logger for errors
	Error *log.Logger
)

func init() {
	Debug = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	Info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lmicroseconds)
	Warn = log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime|log.Lmicroseconds)
	Error = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
}
