package main

import (
	"context"
	"database/sql"
	_ "embed"
	"os"

	"github.com/golang-jwt/jwt/v5"
	db "github.com/haarshmap/go-url/cmd/db/generated"
	"github.com/haarshmap/go-url/internal/server"
	"github.com/joho/godotenv"

	echojwt "github.com/labstack/echo-jwt/v5"
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

	erro := godotenv.Load()
	if erro != nil {
		e.Logger.Error("Failed to initialise env", "error", err)
	}

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

	//grp the sites that are to be protected
	r := e.Group("/dashboard")

	config := echojwt.Config{
		NewClaimsFunc: func(c *echo.Context) jwt.Claims {
			return new(server.JWTCustomClaims)
		},
		SigningKey: []byte(os.Getenv("SIGNING_KEY")),
	}

	r.Use(server.CheckCookie)
	r.Use(echojwt.WithConfig(config))
	r.GET("", h.DashboardHandler)

	if err := e.Start(":" + os.Getenv("PORT")); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:" + os.Getenv("REDISPORT"),
		Password: "",
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		e.Logger.Error("Connecting to redis: %v", "error", err)
	}
}
