# Tab Example

Exercises every tab operation available in go-juggler against a running
camofox-browser instance.

## What it demonstrates

| Step   | Method         | What it does                                        |
| ------ | -------------- | --------------------------------------------------- |
| 1      | `OpenTab`      | Opens a new tab in a session                        |
| 2      | `Navigate`     | Loads a URL (wikipedia.org)                         |
| 3      | `Snapshot`     | Fetches the ARIA accessibility tree with ref IDs    |
| 4      | `Links`        | Lists all links on the page with pagination         |
| 5      | `Stats`        | Returns tab state: URL, visited URLs, ref count     |
| 6      | `Type`         | Fills an input field with text                      |
| 7      | `Press`        | Presses a keyboard key (Enter, Tab, Escape, ...)    |
| 8      | `Scroll`       | Scrolls the page by a pixel amount                  |
| 9      | `Back`         | Navigates back in browser history                   |
| 10     | `Forward`      | Navigates forward in browser history                |
| 11     | `Refresh`      | Reloads the current page                            |
| 12     | `Click`        | Clicks an element by its ARIA ref ID or CSS selector|
| 13     | `Screenshot`   | Captures a PNG screenshot of the page               |
| 14     | `CloseTab`     | Closes the tab                                      |
| 15     | `CloseSession` | Destroys the session and all its tabs               |

## Prerequisites

A running `camofox-browser` container:

```sh
docker run -d -p 9377:9377 --name camoufox camoufox
```

## Run

```sh
go run ./examples/tab -addr http://localhost:9377
```
