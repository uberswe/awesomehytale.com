package handlers

import (
	"net/http"
	"strings"

	"github.com/uberswe/awesomehytale.com/og"
	"github.com/uberswe/awesomehytale.com/parser"
)

// OGHomeHandler generates the OG image for the home page
func (h *Handler) OGHomeHandler(w http.ResponseWriter, r *http.Request) {
	img := og.GenerateHomeImage(h.SiteData.TotalResources, len(h.SiteData.Categories))
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(img)
}

// OGSearchHandler generates the OG image for search
func (h *Handler) OGSearchHandler(w http.ResponseWriter, r *http.Request) {
	img := og.GenerateSearchImage(h.SiteData.TotalResources)
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(img)
}

// OGCategoryHandler generates OG images for category pages
func (h *Handler) OGCategoryHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/og/category/")
	slug := strings.TrimSuffix(path, ".png")

	category := parser.GetCategoryBySlug(h.SiteData, slug)
	if category == nil {
		http.NotFound(w, r)
		return
	}

	img := og.GenerateCategoryImage(category.Name, len(category.Resources), category.Description)
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(img)
}

// OGResourceHandler generates OG images for resource pages
func (h *Handler) OGResourceHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/og/resource/")
	path = strings.TrimSuffix(path, ".png")
	parts := strings.Split(path, "/")

	if len(parts) < 2 {
		http.NotFound(w, r)
		return
	}

	resource := parser.GetResourceBySlug(h.SiteData, parts[0], parts[1])
	if resource == nil {
		http.NotFound(w, r)
		return
	}

	img := og.GenerateResourceImage(resource.Name, resource.Description, resource.Platform, resource.CategoryName)
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(img)
}
