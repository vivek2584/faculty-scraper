package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/vivek2584/faculty-scraper/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	api.RegisterRoutes(mux)

	addr := ":" + port
	fmt.Printf("API server starting on http://localhost%s\n", addr)
	fmt.Println()
	fmt.Println("Endpoints:")
	fmt.Println("  GET /api/campuses                → list campuses")
	fmt.Println("  GET /api/colleges?campus_id=     → list colleges for a campus")
	fmt.Println("  GET /api/departments?college_id= → list departments for a college")
	fmt.Println("  GET /api/department/{id}         → faculty slugs in a department")
	fmt.Println("  GET /api/faculty/{slug}          → full faculty profile JSON")
	fmt.Println("  GET /api/search?name=            → search faculty by name")
	fmt.Println("  GET /faculty/{slug}              → redirect to SRM profile page")
	fmt.Println()

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
