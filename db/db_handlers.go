package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"os"
	"tasks/global"
)

func CreateDBConnection() {
	var err error
	global.DbConn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		global.LogError("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

}
