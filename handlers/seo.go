package handlers

import (
	"fmt"
	"net/http"
	"time"
)

// SitemapHandler generates the XML sitemap
func (h *Handler) SitemapHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")

	now := time.Now().Format("2006-01-02")

	fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>%s/</loc>
    <lastmod>%s</lastmod>
    <changefreq>daily</changefreq>
    <priority>1.0</priority>
  </url>
`, h.BaseURL, now)

	// Add categories
	for _, cat := range h.SiteData.Categories {
		fmt.Fprintf(w, `  <url>
    <loc>%s/category/%s</loc>
    <lastmod>%s</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.8</priority>
  </url>
`, h.BaseURL, cat.Slug, now)

		// Add resources in each category
		for _, res := range cat.Resources {
			fmt.Fprintf(w, `  <url>
    <loc>%s/resource/%s/%s</loc>
    <lastmod>%s</lastmod>
    <changefreq>monthly</changefreq>
    <priority>0.6</priority>
  </url>
`, h.BaseURL, cat.Slug, res.Slug, now)
		}
	}

	fmt.Fprintf(w, `  <url>
    <loc>%s/privacy</loc>
    <lastmod>%s</lastmod>
    <changefreq>yearly</changefreq>
    <priority>0.3</priority>
  </url>
  <url>
    <loc>%s/terms</loc>
    <lastmod>%s</lastmod>
    <changefreq>yearly</changefreq>
    <priority>0.3</priority>
  </url>
</urlset>`, h.BaseURL, now, h.BaseURL, now)
}

// RobotsHandler generates the robots.txt file
func (h *Handler) RobotsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	fmt.Fprintf(w, `User-agent: *
Allow: /
Disallow: /search
Disallow: /go/
Disallow: /static/

Sitemap: %s/sitemap.xml
`, h.BaseURL)
}

// PingSearchEngines notifies search engines about the sitemap
func (h *Handler) PingSearchEngines() {
	sitemapURL := h.BaseURL + "/sitemap.xml"

	// Ping Google
	go func() {
		client := &http.Client{Timeout: 30 * time.Second}
		url := "https://www.google.com/ping?sitemap=" + sitemapURL
		resp, err := client.Get(url)
		if err == nil {
			resp.Body.Close()
			fmt.Printf("Pinged Google sitemap: %s\n", url)
		}
	}()

	// Ping Bing
	go func() {
		client := &http.Client{Timeout: 30 * time.Second}
		url := "https://www.bing.com/ping?sitemap=" + sitemapURL
		resp, err := client.Get(url)
		if err == nil {
			resp.Body.Close()
			fmt.Printf("Pinged Bing sitemap: %s\n", url)
		}
	}()
}
