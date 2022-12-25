package main

import (
	"golang.org/x/exp/slog"
	"os"
	"tasks/global"
)

func SetupLog() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		slog.Error("Failed to open log file", err)
	}
	textHandler := slog.NewTextHandler(os.Stdout)
	global.TextLogger = slog.New(textHandler)
	fileHandler := slog.NewTextHandler(file)
	global.FileLogger = slog.New(fileHandler)
}

func Log(msg string) {
	global.TextLogger.Info(msg)
	global.FileLogger.Info(msg)
}

func LogError(msg string, err error) {
	global.TextLogger.Error(msg, err)
	global.FileLogger.Error(msg, err)
}
