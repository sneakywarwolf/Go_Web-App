package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type Blog struct {
	Title   string
	Content string
	URL     string
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			log.Println("Error executing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	// Serve the static files
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))

	// Start the server
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func getBlogs() []Blog {
	blogFiles, err := filepath.Glob("web/blog/*.html")
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
