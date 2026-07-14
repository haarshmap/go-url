package handlers

import (
	"net/http"

	db "github.com/haarshmap/go-url/cmd/db/generated"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	queries *db.Queries
}

func NewHandler(q *db.Queries) *Handler {
	return &Handler{
		queries: q,
	}
}

// handlers
func (h *Handler) RedirectHandler(c *echo.Context) error {
	id := c.Param("id")

	res, err := h.queries.RedirectHandler(c.Request().Context(), id)
	if err != nil {
		c.Logger().Error("Failed to retrieve url", "error", err)
	}

	return c.Redirect(http.StatusMovedPermanently, res.(string))
}

func (h *Handler) IndexHandler(c *echo.Context) error {
	html := `
		<h1>Submit a new website</h1>
		<form action="/go/submit" method="POST">
		<label for="url">Website URL:</label>
		<input type="text" id="url" name="url">
		<input type="submit" value="Submit">
		</form>`

	return c.HTML(http.StatusOK, html)
}

func (h *Handler) SubmitHandler(c *echo.Context) error {
	return nil
}
