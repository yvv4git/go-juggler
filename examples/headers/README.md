# Headers Example

This example demonstrates how to use custom HTTP headers with the go-juggler library.

## Run

```sh
go run ./examples/headers -addr http://localhost:9377
```

## Description

The example shows:

1. Opening a tab with custom HTTP headers (Authorization, X-Custom-Header, Accept-Language)
2. Navigating to another URL with different headers
3. Taking snapshots to verify the page content

Custom headers are useful for:

- Authentication (Bearer tokens, API keys)
- Custom user agents
- Language preferences
- Request tracking (X-Request-ID)
- Any other HTTP headers needed for your use case
