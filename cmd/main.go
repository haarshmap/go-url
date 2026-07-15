package main

import (
	"context"
	"database/sql"
	_ "embed"
	"os"

	db "github.com/haarshmap/go-url/cmd/db/generated"
	"github.com/haarshmap/go-url/internal/server"

	"github.com/labstack/echo/v5"

	"github.com/labstack/echo/v5/middleware"
	"github.com/redis/go-redis/v9"

	_ "modernc.org/sqlite"
)

//go:embed db/schema.sql
var ddl string

func main() {
	ctx := context.Background()
	e := echo.New()
	var err error

	database, err := sql.Open("sqlite", "data.db")
	if err != nil {
		e.Logger.Error("Failed to initialise database", "error", err)
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		e.Logger.Error("Failed to ping database: %v", "error", err)
	}

	if _, err := database.ExecContext(ctx, ddl); err != nil {
		e.Logger.Error("Failed to create tables: %v", "error", err)
	}

	queries := db.New(database)
	h := server.NewHandler(queries)
	server.RegisterRoutes(e, h)

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	//middleware for custom claims type

	if err := e.Start(":" + os.Getenv("PORT")); err != nil {
		e.Logger.Error("failed to start server: %v", "error", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:" + os.Getenv("RPORT"),
		Password: "",
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		e.Logger.Error("Connecting to redis: %v", "error", err)
	}
}
