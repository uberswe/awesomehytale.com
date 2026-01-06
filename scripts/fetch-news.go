//go:build ignore

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func main() {
	fmt.Println("Fetching news from hytale.com/news...")

	// Fetch the news page
	resp, err := http.Get("https://hytale.com/news")
	if err != nil {
		fmt.Printf("Error fetching news: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		os.Exit(1)
	}

	// Parse news items from the page
	// This is a basic implementation - in production you'd use a proper HTML parser
	content := string(body)

	// Look for news article patterns
	// Note: This is a simplified parser - the actual implementation would need
	// to be adapted based on the real hytale.com/news page structure
	titlePattern := regexp.MustCompile(`<h[12][^>]*>([^<]+)</h[12]>`)
	matches := titlePattern.FindAllStringSubmatch(content, 10)

	for i, match := range matches {
		if len(match) < 2 {
			continue
		}

		title := strings.TrimSpace(match[1])
		if title == "" || len(title) < 10 {
			continue
		}

		// Create slug from title
		slug := strings.ToLower(title)
		slug = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(slug, "-")
		slug = strings.Trim(slug, "-")
		if len(slug) > 50 {
			slug = slug[:50]
		}

		date := time.Now().Format("2006-01-02")
		filename := fmt.Sprintf("%s-%s.md", date, slug)
		filepath := filepath.Join("news", filename)

		// Check if file already exists
		if _, err := os.Stat(filepath); err == nil {
			fmt.Printf("Skipping existing: %s\n", filename)
			continue
		}

		// Create news file
		newsContent := fmt.Sprintf(`---
title: "%s"
source: "hytale.com"
url: "https://hytale.com/news"
date: "%s"
---

News from the official Hytale blog.
`, title, date)

		err := os.WriteFile(filepath, []byte(newsContent), 0644)
		if err != nil {
			fmt.Printf("Error writing %s: %v\n", filename, err)
			continue
		}

		fmt.Printf("Created: %s\n", filename)

		// Limit to 5 new items per run
		if i >= 4 {
			break
		}
	}

	fmt.Println("News fetch complete!")
}
