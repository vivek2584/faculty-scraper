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
	mux.HandleFunc("GET /api/campuses", handleCampuses)
	mux.HandleFunc("GET /api/colleges", handleColleges)
	mux.HandleFunc("GET /api/departments", handleDepartments)
	mux.HandleFunc("GET /api/department/{id}", handleDepartment)
	mux.HandleFunc("GET /api/faculty/{slug}", handleFaculty)
	mux.HandleFunc("GET /api/slug", handleSlug)
	mux.HandleFunc("GET /faculty/{slug}", handleRedirect)
}

// handleCampuses returns the list of SRM campuses.
// GET /api/campuses
func handleCampuses(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, scraper.Campuses)
}

// handleColleges returns colleges for a campus.
// GET /api/colleges?campus_id=78
func handleColleges(w http.ResponseWriter, r *http.Request) {
	campusID := r.URL.Query().Get("campus_id")
	if campusID == "" {
		writeError(w, http.StatusBadRequest, "missing 'campus_id' query parameter")
		return
	}

	colleges, err := scraper.FetchColleges(campusID)
	if err != nil {
		writeError(w, http.StatusBadGateway, "failed to fetch colleges: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, colleges)
}

// handleDepartments returns departments for a college.
// GET /api/departments?college_id=9812
func handleDepartments(w http.ResponseWriter, r *http.Request) {
	collegeID := r.URL.Query().Get("college_id")
	if collegeID == "" {
		writeError(w, http.StatusBadRequest, "missing 'college_id' query parameter")
		return
	}

	depts, err := scraper.FetchDepartments(collegeID)
	if err != nil {
		writeError(w, http.StatusBadGateway, "failed to fetch departments: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, depts)
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

// handleDepartment returns faculty slugs for a department.
// GET /api/department/{id}   (id = WordPress taxonomy ID, e.g. 13519)
func handleDepartment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing department id")
		return
	}

	log.Printf("Fetching faculty slugs for department ID: %s", id)
	slugs, err := scraper.ScrapeDepartmentSlugs(id)
	if err != nil {
		log.Printf("Error fetching department %s: %v", id, err)
		writeError(w, http.StatusBadGateway, "failed to fetch department: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"department_id": id,
		"count":         len(slugs),
		"slugs":         slugs,
	})
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
// GET /api/slug?name={name}
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
