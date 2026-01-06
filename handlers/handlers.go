package handlers

import (
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/uberswe/awesomehytale.com/models"
	"github.com/uberswe/awesomehytale.com/parser"
)

// Handler holds dependencies for HTTP handlers
type Handler struct {
	Templates map[string]*template.Template
	SiteData  *models.SiteData
	BaseURL   string
}

// NewHandler creates a new Handler instance
func NewHandler(templates map[string]*template.Template, siteData *models.SiteData, baseURL string) *Handler {
	return &Handler{
		Templates: templates,
		SiteData:  siteData,
		BaseURL:   baseURL,
	}
}

// HomeHandler handles the home page
func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"Meta": models.Meta{
			Title:        "Awesome Hytale - The Complete Hytale Community Resource Directory",
			Description:  "Discover the best Hytale resources, mods, servers, tools, and communities. The most comprehensive directory for Hytale players, modders, and content creators.",
			CanonicalURL: h.BaseURL + "/",
			OGType:       "website",
			OGImage:      h.BaseURL + "/og/home.png",
			PageType:     "home",
		},
		"SiteData": h.SiteData,
		"BaseURL":  h.BaseURL,
	}

	h.renderTemplate(w, "home.html", data)
}

// CategoryHandler handles category pages
func (h *Handler) CategoryHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/category/")
	slug = strings.TrimSuffix(slug, "/")

	category := parser.GetCategoryBySlug(h.SiteData, slug)
	if category == nil {
		http.NotFound(w, r)
		return
	}

	description := category.Description
	if description == "" {
		description = "Browse " + category.Name + " resources for Hytale. Find the best tools, mods, and community resources."
	}

	data := map[string]interface{}{
		"Meta": models.Meta{
			Title:        category.Name + " - Awesome Hytale",
			Description:  description,
			CanonicalURL: h.BaseURL + "/category/" + slug,
			OGType:       "website",
			OGImage:      h.BaseURL + "/og/category/" + slug + ".png",
			PageType:     "category",
		},
		"Category": category,
		"SiteData": h.SiteData,
		"BaseURL":  h.BaseURL,
	}

	h.renderTemplate(w, "category.html", data)
}

// ResourceHandler handles individual resource pages
func (h *Handler) ResourceHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/resource/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 {
		http.NotFound(w, r)
		return
	}

	categorySlug := parts[0]
	resourceSlug := parts[1]

	resource := parser.GetResourceBySlug(h.SiteData, categorySlug, resourceSlug)
	if resource == nil {
		http.NotFound(w, r)
		return
	}

	description := resource.Description
	if description == "" {
		description = resource.Name + " - A Hytale resource in the " + resource.CategoryName + " category."
	}

	data := map[string]interface{}{
		"Meta": models.Meta{
			Title:        resource.Name + " - Awesome Hytale",
			Description:  truncateDescription(description, 160),
			CanonicalURL: h.BaseURL + "/resource/" + categorySlug + "/" + resourceSlug,
			OGType:       "article",
			OGImage:      h.BaseURL + "/og/resource/" + categorySlug + "/" + resourceSlug + ".png",
			PageType:     "resource",
		},
		"Resource": resource,
		"SiteData": h.SiteData,
		"BaseURL":  h.BaseURL,
	}

	h.renderTemplate(w, "resource.html", data)
}

// SearchHandler handles search requests
func (h *Handler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	var results []models.Resource

	if query != "" {
		results = parser.Search(h.SiteData, query)
	}

	data := map[string]interface{}{
		"Meta": models.Meta{
			Title:        "Search - Awesome Hytale",
			Description:  "Search for Hytale resources, mods, servers, and more.",
			CanonicalURL: h.BaseURL + "/search",
			OGType:       "website",
			OGImage:      h.BaseURL + "/og/search.png",
			PageType:     "search",
		},
		"Query":    query,
		"Results":  results,
		"SiteData": h.SiteData,
		"BaseURL":  h.BaseURL,
	}

	h.renderTemplate(w, "search.html", data)
}

// RedirectHandler handles external redirects with countdown
func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/go/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 {
		http.NotFound(w, r)
		return
	}

	categorySlug := parts[0]
	resourceSlug := parts[1]
	urlIndex := 0

	if len(parts) > 2 {
		if idx, err := strconv.Atoi(parts[2]); err == nil {
			urlIndex = idx
		}
	}

	resource := parser.GetResourceBySlug(h.SiteData, categorySlug, resourceSlug)
	if resource == nil {
		http.NotFound(w, r)
		return
	}

	var targetURL string
	if urlIndex < len(resource.URLs) {
		targetURL = resource.URLs[urlIndex].URL
	} else {
		targetURL = resource.URL
	}

	data := map[string]interface{}{
		"TargetURL":    targetURL,
		"ResourceName": resource.Name,
		"BackURL":      "/resource/" + categorySlug + "/" + resourceSlug,
		"BaseURL":      h.BaseURL,
	}

	h.renderTemplate(w, "redirect.html", data)
}

// PrivacyHandler handles the privacy policy page
func (h *Handler) PrivacyHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Meta": models.Meta{
			Title:        "Privacy Policy - Awesome Hytale",
			Description:  "Privacy policy for Awesome Hytale, the Hytale community resource directory.",
			CanonicalURL: h.BaseURL + "/privacy",
			OGType:       "website",
			OGImage:      h.BaseURL + "/og/home.png",
			PageType:     "legal",
		},
		"SiteData": h.SiteData,
		"BaseURL":  h.BaseURL,
	}

	h.renderTemplate(w, "privacy.html", data)
}

// TermsHandler handles the terms and conditions page
func (h *Handler) TermsHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Meta": models.Meta{
			Title:        "Terms and Conditions - Awesome Hytale",
			Description:  "Terms and conditions for using Awesome Hytale, the Hytale community resource directory.",
			CanonicalURL: h.BaseURL + "/terms",
			OGType:       "website",
			OGImage:      h.BaseURL + "/og/home.png",
			PageType:     "legal",
		},
		"SiteData": h.SiteData,
		"BaseURL":  h.BaseURL,
	}

	h.renderTemplate(w, "terms.html", data)
}

func (h *Handler) renderTemplate(w http.ResponseWriter, name string, data map[string]interface{}) {
	tmpl, ok := h.Templates[name]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		// For redirect template which doesn't use base
		if err.Error() == "html/template: \"base\" is undefined" {
			if err := tmpl.Execute(w, data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func truncateDescription(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// GetEnv gets an environment variable with a default value
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
