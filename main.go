package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/uberswe/awesomehytale.com/handlers"
	"github.com/uberswe/awesomehytale.com/parser"
)

func main() {
	// Load environment configuration
	port := handlers.GetEnv("PORT", "8080")
	baseURL := handlers.GetEnv("BASE_URL", "http://localhost:"+port)
	environment := handlers.GetEnv("ENVIRONMENT", "development")

	fmt.Printf("Starting Awesome Hytale server...\n")
	fmt.Printf("Environment: %s\n", environment)
	fmt.Printf("Base URL: %s\n", baseURL)

	// Parse resources
	siteData, err := parser.ParseResourcesDir("resources")
	if err != nil {
		log.Fatalf("Failed to parse resources: %v", err)
	}
	fmt.Printf("Loaded %d resources across %d categories\n", siteData.TotalResources, len(siteData.Categories))

	// Parse news
	newsItems, err := parser.ParseNewsDir("news")
	if err != nil {
		fmt.Printf("Warning: Could not parse news: %v\n", err)
	} else {
		siteData.NewsItems = newsItems
		fmt.Printf("Loaded %d news items\n", len(newsItems))
	}

	// Load templates
	templates, err := loadTemplates()
	if err != nil {
		log.Fatalf("Failed to load templates: %v", err)
	}
	fmt.Printf("Loaded %d templates\n", len(templates))

	// Create handler
	h := handlers.NewHandler(templates, siteData, baseURL)

	// Set up routes
	http.HandleFunc("/", h.HomeHandler)
	http.HandleFunc("/category/", h.CategoryHandler)
	http.HandleFunc("/resource/", h.ResourceHandler)
	http.HandleFunc("/search", h.SearchHandler)
	http.HandleFunc("/go/", h.RedirectHandler)
	http.HandleFunc("/privacy", h.PrivacyHandler)
	http.HandleFunc("/terms", h.TermsHandler)

	// SEO routes
	http.HandleFunc("/sitemap.xml", h.SitemapHandler)
	http.HandleFunc("/robots.txt", h.RobotsHandler)

	// OG image routes
	http.HandleFunc("/og/home.png", h.OGHomeHandler)
	http.HandleFunc("/og/search.png", h.OGSearchHandler)
	http.HandleFunc("/og/category/", h.OGCategoryHandler)
	http.HandleFunc("/og/resource/", h.OGResourceHandler)

	// Static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Ping search engines in production
	if environment == "production" {
		go h.PingSearchEngines()
	}

	// Start server
	fmt.Printf("Server listening on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func loadTemplates() (map[string]*template.Template, error) {
	templates := make(map[string]*template.Template)

	// Template functions
	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"add": func(a, b int) int {
			return a + b
		},
		"truncate": func(s string, n int) string {
			if len(s) <= n {
				return s
			}
			return s[:n-3] + "..."
		},
	}

	// Load base template with partials
	baseTemplate := template.New("base").Funcs(funcMap)

	// Parse partials
	partials, err := filepath.Glob("templates/partials/*.html")
	if err != nil {
		return nil, err
	}
	for _, partial := range partials {
		_, err := baseTemplate.ParseFiles(partial)
		if err != nil {
			return nil, fmt.Errorf("failed to parse partial %s: %v", partial, err)
		}
	}

	// Parse base template
	_, err = baseTemplate.ParseFiles("templates/base.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse base.html: %v", err)
	}

	// Load page templates
	pageTemplates := []string{
		"home.html",
		"category.html",
		"resource.html",
		"search.html",
		"privacy.html",
		"terms.html",
	}

	for _, page := range pageTemplates {
		// Clone base template
		pageTemplate, err := baseTemplate.Clone()
		if err != nil {
			return nil, fmt.Errorf("failed to clone base template: %v", err)
		}

		// Parse page template
		pagePath := filepath.Join("templates", page)
		_, err = pageTemplate.ParseFiles(pagePath)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %v", pagePath, err)
		}

		templates[page] = pageTemplate
	}

	// Load redirect template separately (doesn't use base)
	redirectTemplate := template.New("redirect").Funcs(funcMap)
	_, err = redirectTemplate.ParseFiles("templates/redirect.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse redirect.html: %v", err)
	}
	templates["redirect.html"] = redirectTemplate

	return templates, nil
}

// For development: reload templates on each request
func init() {
	if os.Getenv("ENVIRONMENT") == "development" {
		fmt.Println("Development mode: templates will be reloaded on each request")
	}
}
