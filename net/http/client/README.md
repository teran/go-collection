# httpclient

HTTP Client for Go to ease interaction with web services, especially APIs

## Usage example

<!-- markdownlint-disable MD010 -->
```go
package main

import (
	"fmt"

	"github.com/pkg/errors"
	ghc "github.com/teran/go-http-client"
)

// FIXME: Change this to your test JSON URL
const testURL = "http://........."

func main() {
	resp := map[string]string{}
	errResp := map[string]string{}
	statusCode, err := ghc.New().
		Base(testURL).
		Get("/json").
		DoJSON(s.ctx, nil, &resp, &errResp)
	if err != nil {
		panic(errors.Wrap(err, fmt.Sprintf("%#v", errResp)))
	}

	fmt.Printf("%#v", resp)
}

```
<!-- markdownlint-enable MD010 -->
