package api

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/vivek2584/faculty-scraper/scraper"
)

var nonAlphaNum = regexp.MustCompile(`[^a-z0-9]+`)

const srmFacultyBase = "https://www.srmist.edu.in/faculty/"

// RegisterRoutes sets up all API routes on the given mux.
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/faculty/{slug}", handleFaculty)
	mux.HandleFunc("GET /api/slug", handleSlug)
	mux.HandleFunc("GET /faculty/{slug}", handleRedirect)
}

// handleFaculty scrapes a single faculty profile and returns full JSON.
// GET /api/faculty/{slug}
func handleFaculty(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		writeError(w, http.StatusBadRequest, "missing faculty slug")
		return
	}

	log.Printf("Scraping faculty profile: %s", slug)
	faculty, err := scraper.ScrapeProfile(slug)
	if err != nil {
		log.Printf("Error scraping %s: %v", slug, err)
		writeError(w, http.StatusBadGateway, "failed to scrape faculty profile: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, faculty)
}

// handleRedirect redirects to the faculty's profile page on the SRM website.
// GET /faculty/{slug}
func handleRedirect(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.Error(w, "missing faculty slug", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, srmFacultyBase+slug+"/", http.StatusTemporaryRedirect)
}

// handleSlug converts a faculty name to a URL-friendly slug.
// GET /api/slug?name=Dr. Ganapathy Sankar U
func handleSlug(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		writeError(w, http.StatusBadRequest, "missing 'name' query parameter")
		return
	}

	slug := strings.ToLower(strings.TrimSpace(name))
	slug = nonAlphaNum.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	writeJSON(w, http.StatusOK, map[string]string{
		"name": name,
		"slug": slug,
	})
}

// writeJSON marshals v as JSON and writes it to the response.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
