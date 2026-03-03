package scraper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/vivek2584/faculty-scraper/models"
)

const (
	userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36"
	baseURL   = "https://www.srmist.edu.in"
)

// ScrapeProfile fetches a single faculty profile by slug.
// slug example: "dr-ganapathy-sankar-u"
func ScrapeProfile(slug string) (*models.Faculty, error) {
	profileURL := baseURL + "/faculty/" + slug + "/"
	var faculty *models.Faculty
	var scrapeErr error

	c := colly.NewCollector(
		colly.AllowedDomains("www.srmist.edu.in"),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", userAgent)
	})

	c.OnError(func(r *colly.Response, err error) {
		scrapeErr = fmt.Errorf("failed to fetch %s: %w", r.Request.URL, err)
	})

	c.OnHTML("html", func(e *colly.HTMLElement) {
		f := parseProfile(e)
		if f.Name != "" {
			faculty = &f
		}
	})

	if err := c.Visit(profileURL); err != nil {
		return nil, fmt.Errorf("failed to visit %s: %w", profileURL, err)
	}
	c.Wait()

	if scrapeErr != nil {
		return nil, scrapeErr
	}
	if faculty == nil {
		return nil, fmt.Errorf("no faculty data found at %s", profileURL)
	}
	return faculty, nil
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
