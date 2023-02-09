package json

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Available(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()
		if query.Get("url") != "example.com" {
			http.Error(rw, fmt.Sprintf("invalid URL: %s", query.Get("url")), http.StatusBadRequest)
			return
		}

		switch query.Get("timestamp") {
		case "20060101":
			file, err := os.Open("./fixtures/timestamp.json")
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			defer func() { _ = file.Close() }()

			_, _ = io.Copy(rw, file)

		case "":
			file, err := os.Open("./fixtures/simple.json")
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			defer func() { _ = file.Close() }()

			_, _ = io.Copy(rw, file)

		default:
			http.Error(rw, fmt.Sprintf("invalid timestamp: %s", query.Get("timestamp")), http.StatusBadRequest)
			return
		}
	}))

	client := New()
	client.baseURL, _ = url.Parse(server.URL)
	client.httpClient = server.Client()

	testCases := []struct {
		desc      string
		host      string
		timestamp string
		expected  *APIResponse
	}{
		{
			desc: "simple",
			host: "example.com",
			expected: &APIResponse{
				URL: "example.com",
				ArchivedSnapshots: Snapshot{
					Closest: Item{
						Status:    "200",
						Available: true,
						URL:       "http://web.archive.org/web/20211231173843/https://www.example.com/",
						Timestamp: "20211231173843",
					},
				},
			},
		},
		{
			desc:      "with timestamp",
			host:      "example.com",
			timestamp: "20060101",
			expected: &APIResponse{
				URL: "example.com",
				ArchivedSnapshots: Snapshot{
					Closest: Item{
						Status:    "200",
						Available: true,
						URL:       "http://web.archive.org/web/20060101213916/http://example.com:80/",
						Timestamp: "20060101213916",
					},
				},
				Timestamp: "20060101",
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			apiResponse, err := client.Available(context.Background(), test.host, test.timestamp)
			require.NoError(t, err)

			assert.Equal(t, test.expected, apiResponse)
		})
	}
}
