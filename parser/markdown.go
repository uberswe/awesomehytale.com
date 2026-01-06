package parser

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/uberswe/awesomehytale.com/models"
)

// Slugify converts a string to a URL-friendly slug
func Slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " & ", "-")
	s = strings.ReplaceAll(s, "&", "-")
	s = strings.ReplaceAll(s, " ", "-")
	// Remove special characters
	reg := regexp.MustCompile(`[^a-z0-9-]`)
	s = reg.ReplaceAllString(s, "")
	// Remove multiple dashes
	reg = regexp.MustCompile(`-+`)
	s = reg.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

// ParseResourcesDir parses all resources from the resources directory
func ParseResourcesDir(resourcesDir string) (*models.SiteData, error) {
	siteData := &models.SiteData{
		Categories: []models.Category{},
	}

	// Walk the resources directory
	entries, err := os.ReadDir(resourcesDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		categoryPath := filepath.Join(resourcesDir, entry.Name())
		category := models.Category{
			Name: formatCategoryName(entry.Name()),
			Slug: entry.Name(),
		}

		// Check for category description
		categoryDescPath := filepath.Join(categoryPath, "_category.md")
		if _, err := os.Stat(categoryDescPath); err == nil {
			desc, _ := parseCategoryDescription(categoryDescPath)
			category.Description = desc
		}

		// Parse resources in this category
		resourceFiles, err := os.ReadDir(categoryPath)
		if err != nil {
			continue
		}

		for _, rf := range resourceFiles {
			if rf.IsDir() || rf.Name() == "_category.md" || !strings.HasSuffix(rf.Name(), ".md") {
				continue
			}

			resourcePath := filepath.Join(categoryPath, rf.Name())
			resource, err := parseResourceFile(resourcePath, category.Name, category.Slug)
			if err != nil {
				continue
			}

			category.Resources = append(category.Resources, *resource)
			siteData.TotalResources++
		}

		siteData.Categories = append(siteData.Categories, category)
	}

	// Sort categories by name
	sort.Slice(siteData.Categories, func(i, j int) bool {
		return siteData.Categories[i].Name < siteData.Categories[j].Name
	})

	return siteData, nil
}

// formatCategoryName converts directory name to display name
func formatCategoryName(dirName string) string {
	words := strings.Split(dirName, "-")
	for i, word := range words {
		if word == "and" || word == "or" {
			continue
		}
		words[i] = strings.Title(word)
	}
	return strings.Join(words, " ")
}

// parseCategoryDescription reads the description from _category.md
func parseCategoryDescription(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inFrontmatter := false
	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" {
			if inFrontmatter {
				break
			}
			inFrontmatter = true
			continue
		}
		if inFrontmatter && strings.HasPrefix(line, "description:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "description:")), nil
		}
	}
	return "", nil
}

// parseResourceFile parses a single resource markdown file
func parseResourceFile(path string, categoryName, categorySlug string) (*models.Resource, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	resource := &models.Resource{
		CategoryName: categoryName,
		CategorySlug: categorySlug,
		Price:        "Free",
		Audience:     "All",
		Platform:     "Web",
	}

	scanner := bufio.NewScanner(file)
	var section string
	var description strings.Builder

	urlPattern := regexp.MustCompile(`\*\*Website:\*\*\s*\[([^\]]*)\]\(([^)]+)\)`)
	platformPattern := regexp.MustCompile(`\*\*Platform\*\*\s*\|\s*(.+)`)
	audiencePattern := regexp.MustCompile(`\*\*Audience\*\*\s*\|\s*(.+)`)
	pricePattern := regexp.MustCompile(`\*\*Price\*\*\s*\|\s*(.+)`)
	categoryPattern := regexp.MustCompile(`\*\*Category:\*\*\s*(.+)`)

	for scanner.Scan() {
		line := scanner.Text()

		// Parse title (# heading)
		if strings.HasPrefix(line, "# ") && resource.Name == "" {
			resource.Name = strings.TrimPrefix(line, "# ")
			resource.Slug = Slugify(resource.Name)
			continue
		}

		// Parse URLs
		if matches := urlPattern.FindStringSubmatch(line); matches != nil {
			link := models.ResourceLink{
				Label: matches[1],
				URL:   matches[2],
			}
			resource.URLs = append(resource.URLs, link)
			if resource.URL == "" {
				resource.URL = matches[2]
			}
			continue
		}

		// Parse category line for subcategory
		if matches := categoryPattern.FindStringSubmatch(line); matches != nil {
			parts := strings.Split(matches[1], " > ")
			if len(parts) > 1 {
				resource.SubcategoryName = strings.TrimSpace(parts[1])
				resource.SubcategorySlug = Slugify(resource.SubcategoryName)
			}
			continue
		}

		// Parse details table
		if matches := platformPattern.FindStringSubmatch(line); matches != nil {
			resource.Platform = strings.TrimSpace(matches[1])
			continue
		}
		if matches := audiencePattern.FindStringSubmatch(line); matches != nil {
			resource.Audience = strings.TrimSpace(matches[1])
			continue
		}
		if matches := pricePattern.FindStringSubmatch(line); matches != nil {
			resource.Price = strings.TrimSpace(matches[1])
			continue
		}

		// Track sections
		if strings.HasPrefix(line, "## ") {
			section = strings.TrimPrefix(line, "## ")
			continue
		}

		// Collect description from Overview section
		if section == "Overview" && line != "" && !strings.HasPrefix(line, "---") {
			description.WriteString(line + " ")
		}
	}

	resource.Description = strings.TrimSpace(description.String())
	return resource, nil
}

// ParseNewsDir parses all news items from the news directory
func ParseNewsDir(newsDir string) ([]models.NewsItem, error) {
	var newsItems []models.NewsItem

	entries, err := os.ReadDir(newsDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") || entry.Name() == ".gitkeep" {
			continue
		}

		newsPath := filepath.Join(newsDir, entry.Name())
		news, err := parseNewsFile(newsPath)
		if err != nil {
			continue
		}

		newsItems = append(newsItems, *news)
	}

	// Sort by date descending
	sort.Slice(newsItems, func(i, j int) bool {
		return newsItems[i].Date.After(newsItems[j].Date)
	})

	return newsItems, nil
}

// parseNewsFile parses a single news markdown file
func parseNewsFile(path string) (*models.NewsItem, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	news := &models.NewsItem{}
	scanner := bufio.NewScanner(file)
	inFrontmatter := false
	var excerpt strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" {
			if inFrontmatter {
				inFrontmatter = false
				continue
			}
			inFrontmatter = true
			continue
		}

		if inFrontmatter {
			if strings.HasPrefix(line, "title:") {
				news.Title = strings.Trim(strings.TrimPrefix(line, "title:"), " \"")
			} else if strings.HasPrefix(line, "source:") {
				news.Source = strings.Trim(strings.TrimPrefix(line, "source:"), " \"")
			} else if strings.HasPrefix(line, "url:") {
				news.URL = strings.Trim(strings.TrimPrefix(line, "url:"), " \"")
			} else if strings.HasPrefix(line, "date:") {
				dateStr := strings.Trim(strings.TrimPrefix(line, "date:"), " \"")
				if t, err := time.Parse("2006-01-02", dateStr); err == nil {
					news.Date = t
				}
			}
		} else if line != "" {
			excerpt.WriteString(line + " ")
		}
	}

	news.Excerpt = strings.TrimSpace(excerpt.String())
	return news, nil
}

// Search searches resources by query
func Search(siteData *models.SiteData, query string) []models.Resource {
	var results []models.Resource
	query = strings.ToLower(query)

	for _, cat := range siteData.Categories {
		for _, res := range cat.Resources {
			if strings.Contains(strings.ToLower(res.Name), query) ||
				strings.Contains(strings.ToLower(res.Description), query) ||
				strings.Contains(strings.ToLower(res.Platform), query) ||
				strings.Contains(strings.ToLower(res.CategoryName), query) {
				results = append(results, res)
			}
		}
	}

	return results
}

// GetResourceBySlug finds a resource by category and resource slug
func GetResourceBySlug(siteData *models.SiteData, categorySlug, resourceSlug string) *models.Resource {
	for _, cat := range siteData.Categories {
		if cat.Slug == categorySlug {
			for _, res := range cat.Resources {
				if res.Slug == resourceSlug {
					return &res
				}
			}
		}
	}
	return nil
}

// GetCategoryBySlug finds a category by slug
func GetCategoryBySlug(siteData *models.SiteData, slug string) *models.Category {
	for _, cat := range siteData.Categories {
		if cat.Slug == slug {
			return &cat
		}
	}
	return nil
}
