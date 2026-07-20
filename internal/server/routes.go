package server

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/haarshmap/go-url/templates"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

var RegisterRoutes = func(e *echo.Echo, h *Handler) {

	config := echojwt.Config{
		NewClaimsFunc: func(c *echo.Context) jwt.Claims {
			return new(JWTCustomClaims)
		},
		SigningKey: []byte(os.Getenv("SIGNING_KEY")),
	}

	e.GET("/", func(c *echo.Context) error {
		var param = c.QueryParam("needs_auth")
		var logged_in = true
		_, err := c.Cookie("access_token")
		if err != nil {
			logged_in = false
		}
		return templates.Home("url-shortie", param, logged_in).Render(c.Request().Context(), c.Response())
	})
	e.GET("/login", func(c *echo.Context) error {
		var param = c.QueryParam("acc_created")
		return templates.Login("url-shortie", param).Render(c.Request().Context(), c.Response())
	}, DenyIfLoggedIn)
	e.GET("/register", func(c *echo.Context) error {
		return templates.Register("url-shortie").Render(c.Request().Context(), c.Response())
	}, DenyIfLoggedIn)

	protected := e.Group("/dashboard")
	protected.Use(CheckCookie)
	// protected.get ==> /dashboard
	protected.GET("", func(c *echo.Context) error {
		return templates.Dashboard("url-shortie").Render(c.Request().Context(), c.Response())
	}, CheckCookie)
	protected.Use(echojwt.WithConfig(config))

	//users routes
	e.POST("/register", h.RegisterHandler, DenyIfLoggedIn)
	e.POST("/login", h.LoginHandler, DenyIfLoggedIn)
	e.POST("/logout", h.LogoutHandler, CheckCookie)
	protected.POST("/dashboard", h.DashboardHandler)

	//links routes
	protected.POST("/dashboard/create", h.CreateLink)
}
