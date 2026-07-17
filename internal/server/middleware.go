package server

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

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
