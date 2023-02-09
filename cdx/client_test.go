package cdx

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Do(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()

		params := map[string]string{
			"limit":         "10",
			"output":        "json",
			"showDupeCount": "true",
			"showResumeKey": "true",
			"showSkipCount": "true",
			"url":           "*.example.com",
		}

		for key, value := range params {
			if query.Get(key) != value {
				http.Error(rw, fmt.Sprintf("invalid %s: %s", key, query.Get(key)), http.StatusBadRequest)
				return
			}
		}

		data, err := os.ReadFile("./fixtures/do.json")
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		_, _ = rw.Write(bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n")))
	}))

	client := New()
	client.baseURL, _ = url.Parse(server.URL)
	client.httpClient = server.Client()

	domain := "*.example.com"

	opts := &APIOptions{
		Output:            "json",
		Limit:             10,
		ShowResumeKey:     true,
		ResumeKey:         "",
		ShowDupeCount:     true,
		ShowSkipCount:     true,
		LastSkipTimestamp: false,
	}

	body, err := client.Do(context.Background(), domain, opts)
	require.NoError(t, err)

	expected := `[["urlkey","timestamp","original","mimetype","statuscode","digest","length","dupecount"],
["com,example)/", "20020120142510", "http://example.com:80/", "text/html", "200", "HT2DYGA5UKZCPBSFVCV3JOBXGW2G5UUA", "1792", "0"],
["com,example)/", "20020328012821", "http://www.example.com:80/", "text/html", "200", "UY3I2DT2AMWAY6DECFCFYMT5ZOTFHUCH", "481", "0"],
["com,example)/", "20020524041628", "http://www.example.com:80/", "text/html", "200", "UY3I2DT2AMWAY6DECFCFYMT5ZOTFHUCH", "481", "1"],
["com,example)/", "20020528114741", "http://www.example.com:80/", "text/html", "200", "UY3I2DT2AMWAY6DECFCFYMT5ZOTFHUCH", "481", "2"],
["com,example)/", "20020529173502", "http://www.example.com:80/", "text/html", "200", "UY3I2DT2AMWAY6DECFCFYMT5ZOTFHUCH", "482", "3"],
["com,example)/", "20020604040806", "http://example.com:80/", "text/html", "200", "UY3I2DT2AMWAY6DECFCFYMT5ZOTFHUCH", "477", "4"],
["com,example)/", "20020604050644", "http://example.com:80/", "text/html", "200", "UY3I2DT2AMWAY6DECFCFYMT5ZOTFHUCH", "476", "5"],
["com,example)/", "20020722232628", "http://example.com:80/", "text/html", "200", "UY3I2DT2AMWAY6DECFCFYMT5ZOTFHUCH", "451", "6"],
["com,example)/", "20020801235910", "http://www.example.com:80/", "text/html", "200", "UY3I2DT2AMWAY6DECFCFYMT5ZOTFHUCH", "459", "7"],
["com,example)/", "20020803080544", "http://www.example.com:80/", "text/html", "200", "UY3I2DT2AMWAY6DECFCFYMT5ZOTFHUCH", "458", "8"],
[],
["com%2Cexample%29%2F+20020803080545"]]
`

	assert.Equal(t, expected, string(body))
}

func TestParseJSON(t *testing.T) {
	data, err := os.ReadFile("./fixtures/do.json")
	require.NoError(t, err)

	items, err := ParseJSON(data)
	require.NoError(t, err)

	expected := []Item{
		{
			URLKey:     "com,example)/",
			Timestamp:  "20020120142510",
			Original:   "http://example.com:80/",
			MimeType:   "text/html",
			StatusCode: 200,
			Digest:     "HT2DYGA5UKZCPBSFVCV3JOBXGW2G5UUA",
			Length:     1792,
			DupeCount:  0,
		},
		{
			URLKey:     "com,example)/",
			Timestamp:  "20020328012821",
			Original:   "http://www.example.com:80/",
			MimeType:   "text/html",
			StatusCode: 200,
			Digest:     "UY3I2DT2AMWAY6DECFCFYMT5ZOTFHUCH",
			Length:     481,
			DupeCount:  0,
		},
		{
			URLKey:     "com,example)/",
			Timestamp:  "20020524041628",
			Original:   "http://www.example.com:80/",
			MimeType:   "text/html",
			StatusCode: 200,
			Digest:     "UY3I2DT2AMWAY6DECFCFYMT5ZOTFHUCH",
			Length:     481,
			DupeCount:  1,
		},
		{
			URLKey:     "com,example)/",
			Timestamp:  "20020528114741",
			Original:   "http://www.example.com:80/",
			MimeType:   "text/html",
			StatusCode: 200,
			Digest:     "UY3I2DT2AMWAY6DECFCFYMT5ZOTFHUCH",
			Length:     481,
			DupeCount:  2,
		},
		{
			URLKey:     "com,example)/",
			Timestamp:  "20020529173502",
			Original:   "http://www.example.com:80/",
			MimeType:   "text/html",
			StatusCode: 200,
			Digest:     "UY3I2DT2AMWAY6DECFCFYMT5ZOTFHUCH",
			Length:     482,
			DupeCount:  3,
		},
		{
			URLKey:     "com,example)/",
			Timestamp:  "20020604040806",
			Original:   "http://example.com:80/",
			MimeType:   "text/html",
			StatusCode: 200,
			Digest:     "UY3I2DT2AMWAY6DECFCFYMT5ZOTFHUCH",
			Length:     477,
			DupeCount:  4,
		},
		{
			URLKey:     "com,example)/",
			Timestamp:  "20020604050644",
			Original:   "http://example.com:80/",
			MimeType:   "text/html",
			StatusCode: 200,
			Digest:     "UY3I2DT2AMWAY6DECFCFYMT5ZOTFHUCH",
			Length:     476,
			DupeCount:  5,
		},
		{
			URLKey:     "com,example)/",
			Timestamp:  "20020722232628",
			Original:   "http://example.com:80/",
			MimeType:   "text/html",
			StatusCode: 200,
			Digest:     "UY3I2DT2AMWAY6DECFCFYMT5ZOTFHUCH",
			Length:     451,
			DupeCount:  6,
		},
		{
			URLKey:     "com,example)/",
			Timestamp:  "20020801235910",
			Original:   "http://www.example.com:80/",
			MimeType:   "text/html",
			StatusCode: 200,
			Digest:     "UY3I2DT2AMWAY6DECFCFYMT5ZOTFHUCH",
			Length:     459,
			DupeCount:  7,
		},
		{
			URLKey:     "com,example)/",
			Timestamp:  "20020803080544",
			Original:   "http://www.example.com:80/",
			MimeType:   "text/html",
			StatusCode: 200,
			Digest:     "UY3I2DT2AMWAY6DECFCFYMT5ZOTFHUCH",
			Length:     458,
			DupeCount:  8,
		},
	}

	assert.Equal(t, expected, items)
}
