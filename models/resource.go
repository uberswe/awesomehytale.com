package models

import "time"

// ResourceLink represents a single URL for a resource
type ResourceLink struct {
	Label string
	URL   string
}

// Resource represents a single Hytale community resource
type Resource struct {
	Name            string
	Slug            string
	URL             string // Primary URL
	Description     string
	URLs            []ResourceLink // All URLs including primary
	Platform        string
	Audience        string
	Price           string
	CategorySlug    string
	SubcategorySlug string
	CategoryName    string
	SubcategoryName string
}

// Subcategory represents a subcategory within a category
type Subcategory struct {
	Name      string
	Slug      string
	Resources []Resource
}

// Category represents a top-level category of resources
type Category struct {
	Name          string
	Slug          string
	Description   string
	Subcategories []Subcategory
	Resources     []Resource // Resources directly in category (no subcategory)
}

// NewsItem represents a news article fetched from official sources
type NewsItem struct {
	Title   string
	Source  string
	URL     string
	Excerpt string
	Date    time.Time
}

// SiteData holds all parsed data for the site
type SiteData struct {
	Categories     []Category
	TotalResources int
	NewsItems      []NewsItem
}

// Meta holds page metadata for templates
type Meta struct {
	Title        string
	Description  string
	CanonicalURL string
	OGType       string
	OGImage      string
	PageType     string // home, category, resource, search, legal
}
