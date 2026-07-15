package server

import (
	"net/http"

	db "github.com/haarshmap/go-url/cmd/db/generated"
	"github.com/labstack/echo/v5"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	queries *db.Queries
}

func NewHandler(q *db.Queries) *Handler {
	return &Handler{
		queries: q,
	}
}

func (h *Handler) RegisterHandler(c *echo.Context) error {
	var req RegisterRequest

	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error while binding")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error while creating hash")
	}

	params := db.CreateUserParams{
		Username:     req.Username,
		HashPassword: string(hash),
		Email:        req.Email,
	}

	user, err := h.queries.CreateUser(c.Request().Context(), params)

	return c.JSON(http.StatusCreated, user)
}

// func (h *Handler) LoginHandler(c *echo.Context) error {
// 	var req RegisterRequest

// 	username := c.FormValue(req.Username)
// 	password := c.FormValue(req.Password)
// 	email := c.Param(req.Email)
// 	user, err := h.queries.GetUserByEmail(c.Request().Context(), email)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, "user not found")
// 	}
// 	passcheck := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(password))
// 	if passcheck != nil || username != user.Username {
// 		return echo.ErrUnauthorized
// 	}

// 	claims := &jwtCustomClaims{
// 		req.Username,
// 		req.Email,
// 		jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// 	t, err := token.SignedString([]byte("secret"))
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, "failed to create a token")
// 	}

// 	return c.JSON(http.StatusOK, map[string]string{
// 		"token": t,
// 	})
// }

// func (h *Handler) DashboardHandler(c *echo.Context) error {
// 	token, err := echo.ContextGet[*jwt.Token](c, "user")
// 	if err != nil {
// 		return echo.ErrUnauthorized.Wrap(err)
// 	}
// 	claims := token.Claims.(*jwtCustomClaims)
// 	name := claims.Name
// 	return c.String(http.StatusOK, "Welcom "+name)
// }

// func (h *Handler) LogoutHandler(c *echo.Context) error {
// 	return nil
// }
