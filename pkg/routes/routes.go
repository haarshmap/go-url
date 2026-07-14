package routes

import (
	"github.com/haarshmap/go-url/pkg/handlers"
	"github.com/labstack/echo/v5"
)

func RegisterRoutes(e *echo.Echo, h *handlers.Handler) error {
	e.GET("/go/:id", h.RedirectHandler)
	e.GET("/", h.IndexHandler)
	e.POST("/go/submit", h.SubmitHandler)

	return nil
}
