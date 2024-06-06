package main

import (
	"fmt"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {

	// Print the current working directory
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
	}
	fmt.Println("Current working directory:", dir)

	r := gin.Default()
	r.MaxMultipartMemory = 1
	// r.LoadHTMLFiles(
	// 	"../../templates/index.html",
	// "../../templates/contact-successful.html",
	// 	".../../templates/contact-failure.html",
	// )

	// Correctly specify the path to the templates
	templatesDir := "../../templates"

	// Use filepath.Walk to load all HTML files in the templates directory
	if err := filepath.Walk(templatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			// Attempt to load the HTML file
			r.LoadHTMLFiles(path)
		}
		return nil
	}); err != nil {
		fmt.Println("Error walking templates directory:", err)
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Contact Form Endpoint
	r.POST("/contact-send", func(c *gin.Context) {
		c.Request.ParseForm()
		email := c.Request.FormValue("email")
		name := c.Request.FormValue("name")
		message := c.Request.FormValue("message")

		// Parse email
		_, err := mail.ParseAddress(email)
		if err != nil {
			c.HTML(http.StatusOK, "contact-failure.html", gin.H{
				// c.JSON(http.StatusBadRequest, gin.H{
				"email": email,
				"error": "invalid email",
			})
		}

		// Make sure name and message is reasonable
		if len(name) > 200 {
			// c.JSON(http.StatusBadRequest, gin.H{
			c.HTML(http.StatusOK, "contact-failure.html", gin.H{
				"name":  name,
				"error": "invalid name",
			})
		}

		if len(message) > 10000 {
			// c.JSON(http.StatusBadRequest, gin.H{
			c.HTML(http.StatusOK, "contact-failure.html", gin.H{
				"message": message,
				"error":   "message too big",
			})
		}

		c.HTML(http.StatusOK, "contact-successful.html", gin.H{
			"name":  name,
			"email": email,
		})

	})

	r.Static("/templates", "../../templates") // can be accessed via http://localhost:8080/templates/

	r.Run()
}
