package server

import (
	"github.com/haarshmap/go-url/templates"
	"github.com/labstack/echo/v5"
)

var RegisterRoutes = func(e *echo.Echo, h *Handler) {
	e.GET("/", func(c *echo.Context) error {
		return templates.Home("url-shortie").Render(c.Request().Context(), c.Response())
	})
	e.GET("/login", func(c *echo.Context) error {
		var param = c.QueryParam("acc_created")
		return templates.Login("url-shortie", param).Render(c.Request().Context(), c.Response())

	})
	e.GET("/register", func(c *echo.Context) error {
		return templates.Register("url-shortie").Render(c.Request().Context(), c.Response())
	})
	e.POST("/register", h.RegisterHandler)
	e.POST("/login", h.LoginHandler)
	e.POST("/logout", h.LogoutHandler)
	e.GET("/dashboard", h.DashboardHandler)
}
