package server

import (
	"github.com/labstack/echo/v5"
)

var RegisterRoutes = func(e *echo.Echo, h *Handler) {
	e.POST("/register", h.RegisterHandler)
	e.POST("/login", h.LoginHandler)
	e.POST("/logout", h.LogoutHandler)
}
