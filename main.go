package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gorm.io/gorm"
)

// models
type Link struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

var db, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{})

func main() {
	e := echo.New()

	//routes
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	e.GET("/:id", RedirectHandler)
	e.GET("/", IndexHandler)
	e.POST("/submit", SubmitHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

// handlers
func RedirectHandler(c echo.Context) error {
	id := "id"

	link, found := linkMap[id]

	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "Link is not found")
	}

	return c.Redirect(http.StatusMovedPermanently, link.URL)
}

func IndexHandler(c echo.Context) error {
	html := `
		<h1>Submit a new website</h1>
		<form action="/submit" method="POST">
		<label for="url">Website URL:</label>
		<input type="text" id="url" name="url">
		<input type="submit" value="Submit">
		</form>
		<h2>Existing Links </h2>
		<ul>`

	for _, link := range linkMap {
		html += `<li><a href="/` + link.ID + `">` + link.ID + `</a></li>`
	}

	html += `</ul>`

	return c.HTML(http.StatusOK, html)
}

func SubmitHandler(c echo.Context) error {
	url := "url"
	if url == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "URL is required")
	}

	if !(len(url) >= 4 && (url[:4] == "http" || url[:5] == "https")) {
		url = "https://" + url
	}

	id := generateRandomString(8)

	linkMap[id] = Link{ID: id, URL: url}

	return c.Redirect(http.StatusSeeOther, "/")

}

// utils
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	var result []byte

	for i := 0; i < length; i++ {
		index := seededRand.Intn(len(charset))
		result = append(result, charset[index])
	}

	return string(result)
}
