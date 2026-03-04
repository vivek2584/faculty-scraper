# SRM Faculty Info Scraper

Go API to scrape faculty info and serve it as JSON payload.

## Running

```
go run main.go
```

Server starts on port 8080 by default. Set the `PORT` environment variable to change it.

## Endpoints

### Discovery (find department IDs)

Find the department ID needed for the faculty listing.

```
GET /api/campuses                     -- list all SRM campuses
GET /api/colleges?campus_id=78        -- colleges under a campus
GET /api/departments?college_id=9812  -- departments under a college
```

### Faculty

```
GET /api/department/{id}              -- all faculty slugs in a department
GET /api/faculty/{slug}               -- full profile for one faculty member
GET /api/search?name=Alice Nithya     -- search faculty by name (returns matching slugs)
GET /faculty/{slug}                   -- redirect to the SRM profile page
```

## Usage

1. Call `/api/campuses` to get campus IDs.
2. Pick a campus and call `/api/colleges?campus_id=...` to list its colleges.
3. Pick a college and call `/api/departments?college_id=...` to list departments.
4. Call `/api/department/{id}` with a department ID to get all faculty slugs.
5. Call `/api/faculty/{slug}` for any slug to get the full profile JSON.
6. Use `/api/search?name=...` to find a faculty member by name. The search queries SRM's live staff-finder and returns faculty slugs.

## Faculty profile fields

The `/api/faculty/{slug}` endpoint returns:

- name, designation, department, campus
- phone, email
- experience, research interest, courses, education, publications, awards, workshops, work experience
- memberships, responsibilities
- image URL, profile URL

## Dependencies

- [colly](https://github.com/gocolly/colly) — base scraper pkg

## License

MIT
