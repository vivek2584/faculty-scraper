package models

// Faculty holds the scraped data for one faculty member.
type Faculty struct {
	Name               string
	Designation        string
	Department         string
	Phone              string
	Email              string
	Campus             string
	Experience         string
	ResearchInterest   string
	Courses            string
	Education          string
	Publications       string
	Awards             string
	Workshops          string
	WorkExperience     string
	Memberships        string
	Responsibilities   string
	ProfileURL         string
}

// CSVHeaders returns the column headers for CSV export.
func CSVHeaders() []string {
	return []string{
		"Name", "Designation", "Department", "Phone", "Email",
		"Campus", "Experience", "Research Interest", "Courses",
		"Education", "Publications", "Awards",
		"Workshops/Seminars/Conferences", "Work Experience",
		"Memberships", "Responsibilities", "Profile URL",
	}
}

// ToRow converts a Faculty to a CSV row.
func (f *Faculty) ToRow() []string {
	return []string{
		f.Name, f.Designation, f.Department, f.Phone, f.Email,
		f.Campus, f.Experience, f.ResearchInterest, f.Courses,
		f.Education, f.Publications, f.Awards,
		f.Workshops, f.WorkExperience,
		f.Memberships, f.Responsibilities, f.ProfileURL,
	}
}
