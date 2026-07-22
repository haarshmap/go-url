package server

import (
	"context"
	"net/http"
	"os"

	"github.com/go-redis/redis_rate/v10"
	"github.com/labstack/echo/v5"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var limiter *redis_rate.Limiter

func InitRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:" + os.Getenv("REDISPORT"),
		Password: "",
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		echo.New().Logger.Error("Connecting to redis:", "error", err)
	}
	limiter = redis_rate.NewLimiter(rdb)
}

func CheckCookie(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		cookie, err := c.Cookie("access_token")
		if err != nil {
			c.Logger().Error("cookie not found")
			return echo.NewHTTPError(http.StatusBadRequest, "cookie not found")
		}

		c.Request().Header.Set("Authorization", "bearer "+cookie.Value)
		return next(c)
	}
}

func IsLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		_, err := c.Cookie("access_token")
		if err == nil {
			c.Logger().Info("User is logged in")
			c.Redirect(http.StatusMovedPermanently, "/dashboard")
		}

		return next(c)
	}
}

func RateLimiter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		ClientIP := c.RealIP()
		r, err := limiter.Allow(c.Request().Context(), ClientIP, redis_rate.PerMinute(20))
		if err != nil {
			c.Logger().Error("Error: ", "error", err)
			return echo.NewHTTPError(http.StatusBadGateway, "error")
		}
		if r.Allowed == 0 {
			c.Logger().Info("Too many requests at once")
			return echo.NewHTTPError(http.StatusTooManyRequests, "Try again later too many requests")
		}
		return next(c)
	}
}
