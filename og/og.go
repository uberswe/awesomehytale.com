package og

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"sync"

	"github.com/fogleman/gg"
)

const (
	Width  = 1200
	Height = 630
)

// Hytale-inspired colors
var (
	BgDark       = color.RGBA{13, 17, 23, 255}    // #0d1117
	BgMedium     = color.RGBA{22, 27, 34, 255}    // #161b22
	AccentCyan   = color.RGBA{0, 188, 212, 255}   // #00bcd4
	AccentOrange = color.RGBA{255, 152, 0, 255}   // #ff9800
	TextLight    = color.RGBA{230, 237, 243, 255} // #e6edf3
	TextMuted    = color.RGBA{139, 148, 158, 255} // #8b949e
)

var cache = struct {
	sync.RWMutex
	images map[string][]byte
}{images: make(map[string][]byte)}

func getCached(key string) ([]byte, bool) {
	cache.RLock()
	defer cache.RUnlock()
	img, ok := cache.images[key]
	return img, ok
}

func setCache(key string, img []byte) {
	cache.Lock()
	defer cache.Unlock()
	cache.images[key] = img
}

func encodePNG(img image.Image) []byte {
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}

func createContext() *gg.Context {
	dc := gg.NewContext(Width, Height)

	// Draw gradient background
	grad := gg.NewLinearGradient(0, 0, 0, Height)
	grad.AddColorStop(0, BgDark)
	grad.AddColorStop(1, color.RGBA{8, 10, 14, 255})
	dc.SetFillStyle(grad)
	dc.DrawRectangle(0, 0, Width, Height)
	dc.Fill()

	// Draw accent line at bottom
	dc.SetColor(AccentCyan)
	dc.DrawRectangle(0, Height-4, Width, 4)
	dc.Fill()

	return dc
}

func loadFont(dc *gg.Context, size float64) {
	// Try to load system fonts
	fonts := []string{
		"/System/Library/Fonts/Helvetica.ttc",
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
		"/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf",
	}
	for _, font := range fonts {
		if err := dc.LoadFontFace(font, size); err == nil {
			return
		}
	}
}

func drawStatsBox(dc *gg.Context, x, y, w, h float64, value, label string) {
	// Draw box background
	dc.SetColor(BgMedium)
	dc.DrawRoundedRectangle(x, y, w, h, 8)
	dc.Fill()

	// Draw border
	dc.SetColor(AccentCyan)
	dc.SetLineWidth(2)
	dc.DrawRoundedRectangle(x, y, w, h, 8)
	dc.Stroke()

	// Draw value
	loadFont(dc, 36)
	dc.SetColor(TextLight)
	dc.DrawStringAnchored(value, x+w/2, y+h/2-10, 0.5, 0.5)

	// Draw label
	loadFont(dc, 18)
	dc.SetColor(TextMuted)
	dc.DrawStringAnchored(label, x+w/2, y+h/2+20, 0.5, 0.5)
}

// GenerateHomeImage creates the home page OG image
func GenerateHomeImage(totalResources, totalCategories int) []byte {
	cacheKey := fmt.Sprintf("home:%d:%d", totalResources, totalCategories)
	if cached, ok := getCached(cacheKey); ok {
		return cached
	}

	dc := createContext()

	// Title
	loadFont(dc, 56)
	dc.SetColor(TextLight)
	dc.DrawStringAnchored("Awesome Hytale", Width/2, 120, 0.5, 0.5)

	// Subtitle
	loadFont(dc, 28)
	dc.SetColor(AccentCyan)
	dc.DrawStringAnchored("The Complete Hytale Community Resource Directory", Width/2, 180, 0.5, 0.5)

	// Stats boxes
	boxWidth := 200.0
	boxHeight := 100.0
	startX := (Width - (boxWidth*3 + 60)) / 2

	drawStatsBox(dc, startX, 280, boxWidth, boxHeight, fmt.Sprintf("%d", totalResources), "Resources")
	drawStatsBox(dc, startX+boxWidth+30, 280, boxWidth, boxHeight, fmt.Sprintf("%d", totalCategories), "Categories")
	drawStatsBox(dc, startX+(boxWidth+30)*2, 280, boxWidth, boxHeight, "Free", "Always")

	// Branding
	loadFont(dc, 20)
	dc.SetColor(TextMuted)
	dc.DrawStringAnchored("awesomehytale.com", Width/2, Height-50, 0.5, 0.5)

	img := encodePNG(dc.Image())
	setCache(cacheKey, img)
	return img
}

// GenerateSearchImage creates the search page OG image
func GenerateSearchImage(totalResources int) []byte {
	cacheKey := "search"
	if cached, ok := getCached(cacheKey); ok {
		return cached
	}

	dc := createContext()

	// Title
	loadFont(dc, 56)
	dc.SetColor(TextLight)
	dc.DrawStringAnchored("Search Resources", Width/2, 200, 0.5, 0.5)

	// Subtitle
	loadFont(dc, 28)
	dc.SetColor(AccentCyan)
	dc.DrawStringAnchored(fmt.Sprintf("Search through %d+ Hytale resources", totalResources), Width/2, 280, 0.5, 0.5)

	// Branding
	loadFont(dc, 20)
	dc.SetColor(TextMuted)
	dc.DrawStringAnchored("awesomehytale.com", Width/2, Height-50, 0.5, 0.5)

	img := encodePNG(dc.Image())
	setCache(cacheKey, img)
	return img
}

// GenerateCategoryImage creates a category page OG image
func GenerateCategoryImage(name string, resourceCount int, description string) []byte {
	dc := createContext()

	// Category name
	loadFont(dc, 48)
	dc.SetColor(TextLight)
	dc.DrawStringAnchored(name, Width/2, 150, 0.5, 0.5)

	// Resource count
	loadFont(dc, 28)
	dc.SetColor(AccentCyan)
	dc.DrawStringAnchored(fmt.Sprintf("%d Resources", resourceCount), Width/2, 220, 0.5, 0.5)

	// Description (truncated)
	if len(description) > 100 {
		description = description[:97] + "..."
	}
	if description != "" {
		loadFont(dc, 22)
		dc.SetColor(TextMuted)
		dc.DrawStringAnchored(description, Width/2, 300, 0.5, 0.5)
	}

	// Branding
	loadFont(dc, 20)
	dc.SetColor(TextMuted)
	dc.DrawStringAnchored("awesomehytale.com", Width/2, Height-50, 0.5, 0.5)

	return encodePNG(dc.Image())
}

// GenerateResourceImage creates a resource page OG image
func GenerateResourceImage(name, description, platform, category string) []byte {
	dc := createContext()

	// Resource name
	loadFont(dc, 44)
	dc.SetColor(TextLight)
	dc.DrawStringAnchored(name, Width/2, 120, 0.5, 0.5)

	// Category badge
	loadFont(dc, 20)
	dc.SetColor(AccentOrange)
	dc.DrawStringAnchored(category, Width/2, 170, 0.5, 0.5)

	// Description (truncated)
	if len(description) > 150 {
		description = description[:147] + "..."
	}
	if description != "" {
		loadFont(dc, 22)
		dc.SetColor(TextMuted)
		dc.DrawStringAnchored(description, Width/2, 260, 0.5, 0.5)
	}

	// Platform info
	loadFont(dc, 24)
	dc.SetColor(AccentCyan)
	dc.DrawStringAnchored("Platform: "+platform, Width/2, 400, 0.5, 0.5)

	// Branding
	loadFont(dc, 20)
	dc.SetColor(TextMuted)
	dc.DrawStringAnchored("awesomehytale.com", Width/2, Height-50, 0.5, 0.5)

	return encodePNG(dc.Image())
}
