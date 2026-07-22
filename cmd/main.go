package main

import (
	"context"
	"database/sql"
	_ "embed"
	"os"

	db "github.com/haarshmap/go-url/cmd/db/generated"
	"github.com/haarshmap/go-url/internal/server"
	"github.com/joho/godotenv"

	"github.com/labstack/echo/v5"

	"github.com/labstack/echo/v5/middleware"

	_ "modernc.org/sqlite"
)

//go:embed db/schema.sql
var ddl string

func main() {
	ctx := context.Background()
	e := echo.New()
	err := godotenv.Load()
	if err != nil {
		e.Logger.Error("Failed to initialise env", "error", err)
	}

	server.InitRedis()

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

	//middleware
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	if err := e.Start(":" + os.Getenv("PORT")); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}

}
