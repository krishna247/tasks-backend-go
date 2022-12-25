package global

import (
	"github.com/jackc/pgx/v5"
	"golang.org/x/exp/slog"
)

var FileLogger *slog.Logger
var TextLogger *slog.Logger

var DbConn *pgx.Conn
