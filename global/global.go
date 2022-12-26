package global

import (
	"github.com/jackc/pgx/v5"
	"golang.org/x/exp/slog"
	"os"
)

var FileLogger *slog.Logger
var TextLogger *slog.Logger

var DbConn *pgx.Conn

func SetupLog() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		slog.Error("Failed to open log file", err)
	}
	textHandler := slog.NewTextHandler(os.Stdout)
	TextLogger = slog.New(textHandler)
	fileHandler := slog.NewTextHandler(file)
	FileLogger = slog.New(fileHandler)
}

func Log(msg string) {
	TextLogger.Info(msg)
	FileLogger.Info(msg)
}

func LogError(msg string, err error) {
	TextLogger.Error(msg, err)
	FileLogger.Error(msg, err)
}
