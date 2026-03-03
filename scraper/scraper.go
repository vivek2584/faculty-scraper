package scraper

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/vivek2584/faculty-scraper/models"
)

const userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36"

// Scrape fetches faculty data from the given URL.
// It automatically detects whether the URL points to a single faculty profile
// or a department listing page and scrapes accordingly.
func Scrape(startURL string) ([]models.Faculty, error) {
	var (
		faculties []models.Faculty
		mu       sync.Mutex
		errs     []string
	)

	// ── Collector for individual faculty profile pages ──────────────
	profileCollector := colly.NewCollector(
		colly.AllowedDomains("www.srmist.edu.in"),
	)

	profileCollector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 1,
		Delay:       2 * time.Second,
	})

	profileCollector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", userAgent)
		fmt.Printf("  → Scraping profile: %s\n", r.URL.String())
	})

	profileCollector.OnError(func(r *colly.Response, err error) {
		mu.Lock()
		errs = append(errs, fmt.Sprintf("error scraping %s: %v", r.Request.URL, err))
		mu.Unlock()
		fmt.Printf("  ✗ Error scraping %s: %v\n", r.Request.URL, err)
	})

	profileCollector.OnHTML("html", func(e *colly.HTMLElement) {
		f := parseProfile(e)
		if f.Name != "" {
			mu.Lock()
			faculties = append(faculties, f)
			mu.Unlock()
			fmt.Printf("  ✓ Found: %s | %s\n", f.Name, f.Designation)
		}
	})

	// ── Main collector for department/listing pages ─────────────────
	mainCollector := colly.NewCollector(
		colly.AllowedDomains("www.srmist.edu.in"),
	)

	mainCollector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 1,
		Delay:       2 * time.Second,
	})

	mainCollector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", userAgent)
		fmt.Printf("Visiting: %s\n", r.URL.String())
	})

	mainCollector.OnError(func(r *colly.Response, err error) {
		mu.Lock()
		errs = append(errs, fmt.Sprintf("error visiting %s: %v", r.Request.URL, err))
		mu.Unlock()
		fmt.Printf("Error: %s - %v\n", r.Request.URL, err)
	})

	// Discover faculty profile links on department pages.
	mainCollector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		absURL := e.Request.AbsoluteURL(link)

		if strings.Contains(absURL, "/faculty/") && !strings.HasSuffix(absURL, "/faculty/") {
			profileCollector.Visit(absURL)
		}
	})

	// ── Dispatch ────────────────────────────────────────────────────
	if isSingleProfile(startURL) {
		fmt.Println("Mode: Single faculty profile")
		profileCollector.Visit(startURL)
	} else {
		fmt.Println("Mode: Department page → discovering faculty profiles...")
		mainCollector.Visit(startURL)
	}

	mainCollector.Wait()
	profileCollector.Wait()

	if len(errs) > 0 {
		return faculties, fmt.Errorf("encountered %d error(s): %s", len(errs), strings.Join(errs, "; "))
	}

	return faculties, nil
}

// isSingleProfile returns true when the URL points to a specific faculty
// member rather than a department listing.
func isSingleProfile(url string) bool {
	return strings.Contains(url, "/faculty/") && !strings.HasSuffix(url, "/faculty/")
}

// parseProfile extracts Faculty data from a profile page's root HTML element.
func parseProfile(e *colly.HTMLElement) models.Faculty {
	f := models.Faculty{
		ProfileURL: e.Request.URL.String(),
	}

	// Name (from og:title meta tag — cleanest source)
	f.Name = e.ChildAttr("meta[property='og:title']", "content")
	f.Name = strings.TrimSuffix(f.Name, " - SRMIST")
	f.Name = strings.TrimSpace(f.Name)

	// Info list items (designation, department, phone, email)
	var listItems []string
	e.ForEach("div.hide_empty_list_item .elementor-icon-list-items .elementor-icon-list-item .elementor-icon-list-text", func(i int, el *colly.HTMLElement) {
		text := strings.TrimSpace(el.Text)
		if text != "" {
			listItems = append(listItems, text)
		}
	})

	if len(listItems) >= 1 {
		f.Designation = listItems[0]
	}
	if len(listItems) >= 2 {
		f.Department = listItems[1]
	}
	if len(listItems) >= 3 {
		f.Phone = listItems[2]
	}
	if len(listItems) >= 4 {
		f.Email = listItems[3]
	}

	// Campus / College info
	campusText := strings.TrimSpace(e.ChildText(".faculty-cdc"))
	campusText = strings.TrimPrefix(campusText, "CAMPUS:")
	campusText = strings.TrimSpace(campusText)
	campusText = strings.Join(strings.Fields(campusText), " ")
	f.Campus = campusText

	// Experience
	e.ForEach("div[data-widget_type='text-editor.default'] .elementor-widget-container", func(i int, el *colly.HTMLElement) {
		text := strings.TrimSpace(el.Text)
		if strings.HasPrefix(text, "EXPERIENCE") {
			exp := strings.TrimPrefix(text, "EXPERIENCE :")
			exp = strings.TrimPrefix(exp, "EXPERIENCE:")
			f.Experience = strings.TrimSpace(exp)
		}
	})

	// Research Interest
	e.ForEach("div[data-widget_type='text-editor.default'] .elementor-widget-container", func(i int, el *colly.HTMLElement) {
		text := strings.TrimSpace(el.Text)
		if strings.HasPrefix(text, "RESEARCH INTEREST") {
			ri := strings.TrimPrefix(text, "RESEARCH INTEREST :")
			ri = strings.TrimPrefix(ri, "RESEARCH INTEREST:")
			f.ResearchInterest = strings.TrimSpace(ri)
		}
	})

	// Courses
	e.ForEach("div[data-widget_type='text-editor.default'] .elementor-widget-container", func(i int, el *colly.HTMLElement) {
		text := strings.TrimSpace(el.Text)
		if strings.HasPrefix(text, "COURSES") {
			c := strings.TrimPrefix(text, "COURSES :")
			c = strings.TrimPrefix(c, "COURSES:")
			f.Courses = strings.TrimSpace(c)
		}
	})

	// ── Tab content ───────────────────────────────────────────────
	// Tab titles map to content panels via data-tab attribute.
	// We read each tab's inner text, collapsing whitespace.

	tabTitles := make(map[int]string)
	e.ForEach(".elementor-tab-title", func(i int, el *colly.HTMLElement) {
		tab := el.Attr("data-tab")
		if tab != "" {
			n, _ := strconv.Atoi(tab)
			title := strings.ToLower(strings.TrimSpace(el.Text))
			tabTitles[n] = title
		}
	})

	e.ForEach(".elementor-tab-content", func(i int, el *colly.HTMLElement) {
		tab := el.Attr("data-tab")
		if tab == "" {
			return
		}
		n, _ := strconv.Atoi(tab)
		title := tabTitles[n]
		content := cleanText(el.Text)
		if content == "" {
			return
		}

		switch {
		case strings.Contains(title, "education"):
			f.Education = content
		case strings.Contains(title, "publication"):
			f.Publications = content
		case strings.Contains(title, "award"):
			f.Awards = content
		case strings.Contains(title, "workshop") || strings.Contains(title, "seminar") || strings.Contains(title, "conference"):
			f.Workshops = content
		case strings.Contains(title, "work experience"):
			f.WorkExperience = content
		case strings.Contains(title, "membership"):
			f.Memberships = content
		case strings.Contains(title, "responsibilities") || strings.Contains(title, "responsibility"):
			f.Responsibilities = content
		}
	})

	return f
}

// cleanText collapses runs of whitespace into single spaces and trims.
func cleanText(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}
