// Package cdx Wayback CDX Server API.
// https://github.com/internetarchive/wayback/tree/master/wayback-cdx-server
package cdx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

// BaseURL base URL of the API endpoint.
const BaseURL = "http://web.archive.org"

// Client is an CDX API client.
type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
}

// New creates a new Client.
func New() *Client {
	baseURL, _ := url.Parse(BaseURL)

	return &Client{
		httpClient: &http.Client{Timeout: 50 * time.Second},
		baseURL:    baseURL,
	}
}

// Do sends request to the API endpoint.
func (c Client) Do(ctx context.Context, host string, opts *APIOptions) ([]byte, error) {
	endpoint := c.baseURL.JoinPath("cdx", "search", "cdx")

	var values url.Values
	if opts == nil {
		values = endpoint.Query()
	} else {
		var err error
		values, err = query.Values(opts)
		if err != nil {
			return nil, err
		}
	}

	values.Set("url", host)

	endpoint.RawQuery = values.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

// ParseJSON parses the JSON response to a struct.
func ParseJSON(body []byte) ([]Item, error) {
	var d [][]string
	err := json.Unmarshal(body, &d)
	if err != nil {
		return nil, err
	}

	if len(d) == 0 {
		return nil, nil
	}

	// field names
	fieldsMapping := getFieldsMapping(d[0])

	var items []Item
	for i, line := range d {
		if i == 0 {
			continue
		}
		if len(line) == 0 {
			break
		}

		item := Item{}

		value := reflect.ValueOf(&item)
		for i, elt := range line {
			name := fieldsMapping[i]
			if name == "" {
				log.Printf("unknown field: position %d, value: %s", i, elt)
				continue
			}

			field := value.Elem().FieldByName(name)

			switch field.Kind() {
			case reflect.String:
				field.SetString(elt)
			case reflect.Int:
				val, _ := strconv.ParseInt(elt, 10, 64)
				field.SetInt(val)
			case reflect.Struct:
				if field.Type() == reflect.TypeOf(time.Time{}) {
					ts, err := time.Parse("20060102150405", elt)
					if err != nil {
						return nil, err
					}
					field.Set(reflect.ValueOf(ts))
				}
			default:
				return nil, fmt.Errorf("unsupported type %q for field %q position %d value %q", field.Kind(), name, i, elt)
			}
		}

		items = append(items, item)
	}

	return items, nil
}

func getFieldsMapping(titles []string) map[int]string {
	jsonMapping := make(map[string]string)
	itemType := reflect.TypeOf(Item{})
	for j := 0; j < itemType.NumField(); j++ {
		field := itemType.Field(j)
		jsonMapping[strings.ReplaceAll(field.Tag.Get("json"), ",omitempty", "")] = field.Name
	}

	fieldsMapping := make(map[int]string)
	for i, s := range titles {
		fieldsMapping[i] = jsonMapping[s]
	}

	return fieldsMapping
}
