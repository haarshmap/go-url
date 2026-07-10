package main

import (
	"errors"
	"math/rand/v2"
	"net/http"

	"github.com/glebarez/sqlite"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gorm.io/gorm"
)

// models
type Link struct {
	gorm.Model
	SID string `gorm:"uniqueIndex, not null" json:"sid"`
	URL string `gorm:"not null" json:"url"`
}

var db *gorm.DB

// TODO : error handling at main func and database to not be global
func main() {
	e := echo.New()
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		echo.NewHTTPError(http.StatusNotFound, "database not initialised")
	}
	db.AutoMigrate(&Link{})

	//routes
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	e.GET("/go/:id", RedirectHandler)
	e.GET("/", IndexHandler)
	e.POST("/go/submit", SubmitHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

// handlers
func RedirectHandler(c echo.Context) error {
	return nil
}

func IndexHandler(c echo.Context) error {
	html := `
		<h1>Submit a new website</h1>
		<form action="/go/submit" method="POST">
		<label for="url">Website URL:</label>
		<input type="text" id="url" name="url">
		<input type="submit" value="Submit">
		</form>`

	return c.HTML(http.StatusOK, html)
}

func SubmitHandler(c echo.Context) error {
	ID := generateUniqueID(db)
	Url := c.FormValue("url")
	err := db.Create(&Link{
		SID: ID,
		URL: Url,
	}).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "The data was not pushed")
	}
	return c.HTML(http.StatusOK,
		"<p> The shortened code is </p>",
	)
}

// utils
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(length int) string {
	result := make([]byte, length)

	for i := range result {
		result[i] = charset[rand.IntN(len(charset))]
	}

	return string(result)
}

func generateUniqueID(db *gorm.DB) string {
	var link Link
	ID := generateRandomString(8)
	err := db.First(&link, "sid=?", ID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ID
	} else {
		return generateUniqueID(db)
	}

	return ID
}
