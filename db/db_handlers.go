package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
	"tasks/global"
)

func CreateDBConnection() {
	var err error
	global.DbConn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

}
