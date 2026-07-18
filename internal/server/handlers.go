package server

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

	c.Logger().Info("after binding: ",
		"username", req.Username,
		"email", req.Email,
		"password", req.Password,
	)

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
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user failed to be created")
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *Handler) LoginHandler(c *echo.Context) error {
	var req RegisterRequest

	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error while binding")
	}

	c.Logger().Info("after binding: ",
		"email", req.Email,
		"password", req.Password,
	)

	user, err := h.queries.GetUserByEmail(c.Request().Context(), req.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user not found")
	}
	passcheck := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(req.Password))
	if passcheck != nil {
		c.Logger().Info("Password not correct")
		return echo.ErrUnauthorized
	} else {
		c.Logger().Info("User logged in")
	}

	claims := &JWTCustomClaims{
		req.Username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("SIGNING_KEY")))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to sign a token")
	}

	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    t,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, "token is set")
}

func (h *Handler) LogoutHandler(c *echo.Context) error {
	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		MaxAge:   -1,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, "logged out")
}

func (h *Handler) DashboardHandler(c *echo.Context) error {

	return c.JSON(http.StatusFound, "found the cookie and currently in dashboard")
}
