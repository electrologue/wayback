package json

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Available(t *testing.T) {
	client := New()

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
						URL:       "http://web.archive.org/web/20211117020511/https://www.example.com/",
						Timestamp: "20211117020511",
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
