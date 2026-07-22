package server

import (
	"net/http"
	"os"
	"regexp"
	"strconv"
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

	// VALIDATE USER DETAILS
	//for email
	var EmailCheck = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	//for password
	var CapsLetterCheck = regexp.MustCompile(`[A-Z]`)
	var NumCharCheck = regexp.MustCompile(`[0-9]`)
	var SpecialCharCheck = regexp.MustCompile(`[#?!@$%^&*-]`)

	c.Logger().Info("after binding: ",
		"username", req.Username,
		"email", req.Email,
		"password", req.Password,
	)

	if !EmailCheck.MatchString(req.Email) {
		c.Logger().Error("Invalid email address")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Email address")
	}

	if !CapsLetterCheck.MatchString(req.Password) {
		c.Logger().Error("Must Include a Capital Letter")
		return echo.NewHTTPError(http.StatusBadRequest, "Must Include a Capital Letter")
	}

	if !NumCharCheck.MatchString(req.Password) {
		c.Logger().Error("Must Include a Number")
		return echo.NewHTTPError(http.StatusBadRequest, "Must Include a Number")
	}

	if !SpecialCharCheck.MatchString(req.Password) {
		c.Logger().Error("Must Include a Special Character")
		return echo.NewHTTPError(http.StatusBadRequest, "Must Include a Special Character")
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

	_, err = h.queries.CreateUser(c.Request().Context(), params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user failed to be created")
	}

	c.Logger().Info("Account created! You can login.")
	return c.JSON(http.StatusOK, "/login?acc_created=true")
}

func (h *Handler) LoginHandler(c *echo.Context) error {
	var login LoginRequest

	err := c.Bind(&login)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error while binding")
	}

	if login.Email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "email not found")
	}

	c.Logger().Info("after binding: ",
		"email", login.Email,
		"password", login.Password,
	)

	user, err := h.queries.GetUserByEmail(c.Request().Context(), login.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user not found")
	}
	passcheck := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(login.Password))
	if passcheck != nil {
		c.Logger().Info("Password not correct")
		return echo.ErrUnauthorized
	} else {
		c.Logger().Info("User logged in")
	}

	claims := &JWTCustomClaims{
		login.Username,
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

	return c.JSON(http.StatusOK, "/dashboard")
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

	return c.JSON(http.StatusOK, "user is logged out")
}

func (h *Handler) DashboardHandler(c *echo.Context) error {

	return c.JSON(http.StatusFound, "found the cookie and currently in dashboard")
}

func (h *Handler) RegisterLink(c *echo.Context) error {
	var link Links

	c.Logger().Info("Before binding:",
		"link.Short_code:", link.ShortID,
		"link.Original_URL:", link.OrigURL,
		"link.UserID:", link.UserID,
		"link.Expiry:", link.Expiry,
	)

	err := c.Bind(&link)
	c.Logger().Info("After binding:",
		"link.Short_code:", link.ShortID,
		"link.Original_URL:", link.OrigURL,
		"link.UserID:", link.UserID,
		"link.Expiry:", link.Expiry,
	)
	if err != nil {
		c.Logger().Error("Failed to bind links", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to bind")
	}

	link.Expiry = time.Now().Add(24 * time.Hour)

	if link.ShortID == "" {
		link.ShortID = strconv.FormatUint(uint64(FastHash(link.OrigURL)), 10)
	}

	params := db.CreateLinkParams{
		ShortID: link.ShortID,
		OrigUrl: link.OrigURL,
		Expiry:  link.Expiry,
		UserID:  link.UserID,
	}

	_, err = h.queries.CreateLink(c.Request().Context(), params)
	if err != nil {
		c.Logger().Error("Failed to create the user")
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to create a link row")
	}

	return c.JSON(http.StatusOK, "link created")
}
