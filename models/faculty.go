package models

// Option represents an ID/title pair from the SRM staff-finder dropdowns.
type Option struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// Faculty holds the scraped data for one faculty member.
type Faculty struct {
	Name             string `json:"name"`
	Designation      string `json:"designation"`
	Department       string `json:"department"`
	Phone            string `json:"phone,omitempty"`
	Email            string `json:"email,omitempty"`
	Campus           string `json:"campus"`
	Experience       string `json:"experience,omitempty"`
	ResearchInterest string `json:"research_interest,omitempty"`
	Courses          string `json:"courses,omitempty"`
	Education        string `json:"education,omitempty"`
	Publications     string `json:"publications,omitempty"`
	Awards           string `json:"awards,omitempty"`
	Workshops        string `json:"workshops,omitempty"`
	WorkExperience   string `json:"work_experience,omitempty"`
	Memberships      string `json:"memberships,omitempty"`
	Responsibilities string `json:"responsibilities,omitempty"`
	ImageURL         string `json:"image_url,omitempty"`
	ProfileURL       string `json:"profile_url"`
}
