package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/gorm/logger"
)

const (
	LOG_DIR       = "./logs/query_log"
	ERROR_LOG_DIR = "./logs/error_log"
)

var errorLogFile *os.File

func SetupLogger() logger.Interface {
	err := os.MkdirAll(LOG_DIR, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create log directory: %v", err)
	}

	currentMonth := time.Now().Format("January")
	currentMonth = strings.ToLower(currentMonth)
	logFileName := currentMonth + "_query.log"

	logFile, err := os.OpenFile(LOG_DIR+"/"+logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}

	newLogger := logger.New(
		log.New(logFile, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      false,
		},
	)

	return newLogger
}

func SetupErrorLogger() error {
	err := os.MkdirAll(ERROR_LOG_DIR, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create error log directory: %v", err)
	}

	currentMonth := time.Now().Format("January")
	currentMonth = strings.ToLower(currentMonth)
	logFileName := currentMonth + "_error.log"

	errorLogFile, err = os.OpenFile(ERROR_LOG_DIR+"/"+logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open error log file: %v", err)
	}

	log.SetOutput(errorLogFile)
	return nil
}

func LogError(format string, args ...interface{}) {
	if errorLogFile == nil {
		log.Printf("[ERROR] "+format, args...)
		return
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("[%s] [ERROR] "+format, append([]interface{}{timestamp}, args...)...)
}
