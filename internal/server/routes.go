package server

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/haarshmap/go-url/templates"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

var RegisterRoutes = func(e *echo.Echo, h *Handler) {

	e.GET("/", func(c *echo.Context) error {
		var logged_in = true
		_, err := c.Cookie("access_token")
		if err != nil {
			logged_in = false
		}
		return templates.Home("url-shortie", logged_in).Render(c.Request().Context(), c.Response())
	})
	e.GET("/login", func(c *echo.Context) error {
		var param = c.QueryParam("acc_created")
		return templates.Login("url-shortie", param).Render(c.Request().Context(), c.Response())
	}, IsLoggedIn)
	e.GET("/register", func(c *echo.Context) error {
		return templates.Register("url-shortie").Render(c.Request().Context(), c.Response())
	}, IsLoggedIn)

	protected := e.Group("/dashboard")

	config := echojwt.Config{
		NewClaimsFunc: func(c *echo.Context) jwt.Claims {
			return new(JWTCustomClaims)
		},
		SigningKey: []byte(os.Getenv("SIGNING_KEY")),
	}

	protected.Use(CheckCookie)
	protected.GET("", func(c *echo.Context) error {
		var logged_in = true
		cookie, err := c.Cookie("access_token")
		if err != nil {
			c.Logger().Error("Cookie is not found")
			return echo.NewHTTPError(http.StatusNotFound, "cookie not found")
		}

		token, err := jwt.ParseWithClaims(cookie.Value, &JWTCustomClaims{}, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusConflict, "Failed to parse")
			}
			return []byte(os.Getenv("SIGNING_KEY")), nil
		})
		if err != nil || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		claims, ok := token.Claims.(*JWTCustomClaims)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to cast claims")
		}

		return templates.Dashboard("url-shortie", logged_in, claims.UserID, claims.Email).Render(c.Request().Context(), c.Response())
	}, CheckCookie)
	protected.Use(echojwt.WithConfig(config))

	//users routes
	e.POST("/register", h.RegisterHandler, IsLoggedIn, RateLimiter)
	e.POST("/login", h.LoginHandler, IsLoggedIn, RateLimiter)
	e.POST("/logout", h.LogoutHandler, CheckCookie)

	//link routes
	// /dashboard
	protected.POST("/", h.DashboardHandler)
	protected.POST("/create", h.RegisterLink)
	e.GET("/:id", h.RedirectLink) //placeholder to check if the redirect works
}
