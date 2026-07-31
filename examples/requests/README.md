# examples/requests

Demonstrates intercepting companion network requests (CSS, JS, images, trackers)
loaded by a page using the Performance API with polling and deduplication.

## Usage

```bash
# Basic — polls for 15s (default)
go run ./examples/requests https://www.wikipedia.org

# Custom wait time
go run ./examples/requests -wait 20s https://rutube.ru

# All flags
go run ./examples/requests -addr http://localhost:9377 -wait 30s -session my-session https://example.com
```

## Flags

| Flag       | Default                     | Description                                       |
|------------|-----------------------------|---------------------------------------------------|
| `-addr`    | `http://localhost:9377`     | camofox-browser endpoint                          |
| `-url`     | `http://www.wikipedia.org`  | page to load (or pass URL as positional arg)      |
| `-session` | `requests-demo`             | session key                                       |
| `-wait`    | `15s`                       | how long to poll for late requests                |

## Output

Requests in chronological order, including the main document:

```text
opened https://www.wikipedia.org/

polling for 15s ...

total requests: 6

 1. [NAVIGATION  ] https://www.wikipedia.org/
 2. [CSS         ] https://www.wikipedia.org/portal/wikipedia.org/assets/img/sprite-e49fbf32.svg
 3. [SCRIPT      ] https://www.wikipedia.org/portal/wikipedia.org/assets/js/index-34f340e24a.js
 4. [SCRIPT      ] https://www.wikipedia.org/portal/wikipedia.org/assets/js/gt-ie9-507b16b6be.js
 5. [IMG         ] https://www.wikipedia.org/portal/wikipedia.org/assets/img/Wikipedia-logo-v2.png
 6. [IMG         ] https://www.wikipedia.org/static/favicon/wikipedia.ico
```

## Notes

- The first entry is always the main document (the page itself).
- The `-wait` flag polls every 1s to capture late-loading trackers and ads.
- `performance.getEntriesByType("resource")` does not show failed/cancelled/blocked requests or websockets — the count will always be lower than DevTools.
