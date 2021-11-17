package cdx

// MatchType allowed values for matchType parameter.
type MatchType string

// https://github.com/internetarchive/wayback/tree/master/wayback-cdx-server#url-match-scope
//  - `matchType=exact` (default if omitted) will return results matching exactly archive.org/about/
//  - `matchType=prefix` will return results for all results under the path `archive.org/about/`
//  - `matchType=host` will return results from host `archive.org`
//  - `matchType=domain` will return results from host archive.org and all subhosts `*.archive.org`
const (
	MatchTypeExact  MatchType = "exact"
	MatchTypePrefix MatchType = "prefix"
	MatchTypeHost   MatchType = "host"
	MatchTypeDomain MatchType = "domain"
)

// APIOptions the API options.
type APIOptions struct {
	// MatchType https://github.com/internetarchive/wayback/tree/master/wayback-cdx-server#url-match-scope
	MatchType MatchType `url:"matchType,omitempty"`

	// Output https://github.com/internetarchive/wayback/tree/master/wayback-cdx-server#output-format-json
	Output string `url:"output,omitempty"`

	// Fields https://github.com/internetarchive/wayback/tree/master/wayback-cdx-server#field-order
	// `["urlkey","timestamp","original","mimetype","statuscode","digest","length"]`
	Fields []string `url:"fl,comma,omitempty"`

	// https://github.com/internetarchive/wayback/tree/master/wayback-cdx-server#filtering
	From   string `url:"from,omitempty"`
	To     string `url:"to,omitempty"`
	Filter string `url:"filter,omitempty"`

	// Collapse https://github.com/internetarchive/wayback/tree/master/wayback-cdx-server#collapsing
	Collapse string `url:"collapse,omitempty"`

	// Limit https://github.com/internetarchive/wayback/tree/master/wayback-cdx-server#query-result-limits
	Limit      int  `url:"limit,omitempty"`
	FastLatest bool `url:"fastLatest,omitempty"`
	Offset     int  `url:"offset,omitempty"`

	// https://github.com/internetarchive/wayback/blob/master/wayback-cdx-server/README.md#resumption-key
	ShowResumeKey bool   `url:"showResumeKey,omitempty"`
	ResumeKey     string `url:"resumeKey,omitempty"`

	// https://github.com/internetarchive/wayback/blob/master/wayback-cdx-server/README.md#counters
	ShowDupeCount     bool `url:"showDupeCount,omitempty"`
	ShowSkipCount     bool `url:"showSkipCount,omitempty"`
	LastSkipTimestamp bool `url:"lastSkipTimestamp,omitempty"`

	// https://github.com/internetarchive/wayback/blob/master/wayback-cdx-server/README.md#pagination-api
	Page           int  `url:"page,omitempty"`
	ShowNumPages   bool `url:"showNumPages,omitempty"`
	PageSize       int  `url:"pageSize,omitempty"`
	ShowPagedIndex bool `url:"showPagedIndex,omitempty"`
}

// Item a history item.
type Item struct {
	URLKey    string `json:"urlkey,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	// Timestamp  time.Time `json:"timestamp,omitempty"`
	Original   string `json:"original,omitempty"`
	MimeType   string `json:"mimetype,omitempty"`
	StatusCode int    `json:"statuscode,omitempty"`
	Digest     string `json:"digest,omitempty"`
	Length     int    `json:"length,omitempty"`

	DupeCount int `json:"dupecount,omitempty"`
}
