package json

// APIResponse the API response.
type APIResponse struct {
	URL               string   `json:"url"`
	ArchivedSnapshots Snapshot `json:"archived_snapshots"`
	Timestamp         string   `json:"timestamp"`
}

// Snapshot snapshot information.
type Snapshot struct {
	Closest Item `json:"closest"`
}

// Item information.
type Item struct {
	Status    string `json:"status"`
	Available bool   `json:"available"`
	URL       string `json:"url"`
	Timestamp string `json:"timestamp"`
}
