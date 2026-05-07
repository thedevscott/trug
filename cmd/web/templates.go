package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/thedevscott/trug/internal/models"
	"github.com/thedevscott/trug/ui"
)

type templateData struct {
	CurrentYear     int
	Transaction     models.Transaction
	Transactions    []models.Transaction
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
	User            models.User
	Stats           models.TransactionStats
}

// humanDate returns a nicely formatted string representation of a time.Time
// value
func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006")
}

func centsToDollar(v int64) string {
	value := float64(v) / 100
	return fmt.Sprintf("%.2f", value)
}

func shortenTitle(title string) string {
	if len(title) > titleLength() {
		return strings.TrimSpace(title[:titleLength()]) + "..."
	}
	return title
}

func titleLength() int {
	return 14
}

var functions = template.FuncMap{
	"humanDate":     humanDate,
	"centsToDollar": centsToDollar,
	"shortenTitle":  shortenTitle,
	"titleLength":   titleLength,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Extract filename ie: home.tmpl.html
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl.html",
			"html/partials/*.tmpl.html",
			page,
		}
		// register template.FuncMap function with template
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
