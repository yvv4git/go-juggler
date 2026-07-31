# Basic Example

Demonstrates the full tab lifecycle against a running Camoufox instance
exposed via the [camofox-browser](https://github.com/jo-inc/camofox-browser)
HTTP wrapper.

## What it does

1. **Health check** -- verifies the browser is alive (`GET /health`).
2. **Open tab** -- creates a new tab in a session (`POST /tabs/open`).
3. **Navigate** -- loads a URL in the tab (`POST /tabs/:id/navigate`).
4. **Snapshot** -- fetches an ARIA snapshot of the page
   (`GET /tabs/:id/snapshot`). The snapshot is an accessibility tree
   with ref IDs (`[e1]`, `[e2]`, ...) that AI agents and automation
   scripts use to locate elements without CSS selectors.
5. **Close tab** -- tears down the tab (`DELETE /tabs/:id`).
6. **Close session** -- destroys the browser context
   (`DELETE /sessions/:id`).

## Prerequisites

A running `camofox-browser` container:

```sh
docker run -d -p 9377:9377 --name camoufox camoufox
```

## Run

```sh
go run ./examples/basic -addr http://localhost:9377
```

## Sample output

```text
health: map[browserConnected:true engine:camoufox ok:true sessions:0]
tab opened: 60bf9c8c-e6d9-4be0-ba09-738058215af6
navigated to rutube.ru
snapshot (128908 chars):
- img
- banner:
  - button "Закрыть меню навигации" [e1]:
    ...
  - button "Вход" [e9]
- navigation:
  - list:
    - listitem:
      - link "Главная" [e10]:
tab closed
session closed
```
