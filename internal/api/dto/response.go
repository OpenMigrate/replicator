package dto

// App is the response shape for a single app.
type App struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// AppList is the response shape for list apps.
type AppList struct {
	NextCursor string `json:"next_cursor"`
	Items      []App  `json:"items"`
}

// Server is the response shape for a server within an app listing.
type Server struct {
	ID           string `json:"id"`
	Hostname     string `json:"hostname"`
	OS           string `json:"os"`
	Arch         string `json:"arch"`
	NumCPU       int    `json:"num_cpu"`
	TimestampUTC string `json:"timestamp_utc"`
}

// ServerList is the response shape for listing servers in an app.
type ServerList struct {
	Total      int64    `json:"total"`
	NextCursor string   `json:"next_cursor"`
	Items      []Server `json:"items"`
}

// Status is a generic OK/ERR style response.
type Status struct {
	Status string `json:"status"`
}

// StatusCount is used where we also want to return a count.
type StatusCount struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}
