# Tab Example

Exercises every tab operation available in go-juggler against a running
camofox-browser instance.

## What it demonstrates

| Step | Method       | What it does                                        |
|------|--------------|-----------------------------------------------------|
| 1    | `OpenTab`    | Opens a new tab in a session                        |
| 2    | `Snapshot`   | Fetches the ARIA accessibility tree with ref IDs    |
| 3    | `Links`      | Lists all links on the page with pagination         |
| 4    | `Stats`      | Returns tab state: URL, visited URLs, ref count     |
| 5    | `Evaluate`   | Injects a button via JavaScript                     |
| 5    | `Click`      | Clicks the injected button by ref                   |
| 6    | `Navigate`   | Loads a URL (example.org)                           |
| 7    | `Back`       | Navigates back in browser history                   |
| 8    | `Forward`    | Navigates forward in browser history                |
| 9    | `Evaluate`   | Injects a form with input field                     |
| 10   | `Snapshot`   | Fetches updated ARIA tree with new refs             |
| 11   | `Type`       | Fills an input field with text                      |
| 12   | `Press`      | Presses Enter to submit the form                    |
| 13   | `Scroll`     | Scrolls the page by a pixel amount                  |
| 14   | `Refresh`    | Reloads the current page                            |
| 15   | `Screenshot` | Captures a PNG screenshot of the page               |

## Prerequisites

A running `camofox-browser` container:

```sh
docker run -d -p 9377:9377 --name camoufox camoufox
```

## Run

```sh
go run ./examples/tab -addr http://localhost:9377
```
