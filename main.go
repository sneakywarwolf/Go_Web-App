package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type Blog struct {
	Title   string
	Content string
	URL     string
}

func main() {
	router := gin.Default()

	// Set up the template engine
	router.SetHTMLTemplate(template.Must(template.ParseFiles("templates/index.html")))

	// Define the route for the homepage
	router.GET("/", func(c *gin.Context) {
		blogs := getBlogs()

		data := struct {
			Title string
			Blogs []Blog
		}{
			Title: "My Blog",
			Blogs: blogs,
		}

		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			log.Println("Error parsing template:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(c.Writer, data)
		if err != nil {
			log.Println("Error executing template:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	})

	// Serve the blog files
	//router.StaticFS("/web", http.Dir("web")) //this wil server the whole web directory to the server
	router.Static("/web/images", "./web/images")
	router.Static("/web/blog", "./web/blog")

	// Start the server
	err := router.Run(":8888")
	if err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}

func getBlogs() []Blog {
	blogFiles, err := filepath.Glob("./web/blog/*.html")
	if err != nil {
		log.Println("Error retrieving blog files:", err)
		return nil
	}

	blogs := make([]Blog, len(blogFiles))
	for i, file := range blogFiles {
		filename := filepath.Base(file)
		title := strings.TrimSuffix(filename, ".html")
		url := "/web/blog/" + filename

		blogs[i] = Blog{
			Title:   title,
			Content: "",
			URL:     url,
		}
	}

	return blogs
}
