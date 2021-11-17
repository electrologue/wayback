# Wayback Machine API clients

##  Wayback Availability JSON API

https://archive.org/help/wayback_api.php

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/electrologue/wayback/json"
)

func main() {
	client := json.New()

	apiResponse, err := client.Available(context.Background(), "*.example.com", "")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", apiResponse)
}

```

## Wayback CDX Server API

https://github.com/internetarchive/wayback/tree/master/wayback-cdx-server

```go
package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/electrologue/wayback/cdx"
)

func main() {
	client := cdx.New()

	domain := "*.example.com"

	opts := &cdx.APIOptions{
		Output:    "json",
		Fields:    []string{"urlkey", "timestamp", "original", "mimetype", "statuscode", "digest", "length"},
		Limit:     4,
	}

	body, err := client.Do(context.Background(), domain, opts)
	if err != nil {
		log.Fatal(err)
	}

	items, err := cdx.ParseJSON(body)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range items {
		if strings.Contains(item.URLKey, `"`) {
			continue
		}

		fmt.Printf("%#v\n", item)
	}
}
```
