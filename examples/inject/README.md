# examples/inject

Demonstrates injecting JavaScript into a page via `POST /tabs/:tabId/evaluate`.

The script sets `document.title` and reads it back — a verifiable side effect.

**Requires camofox-browser >= 1.4.0**

```bash
go run .
```

Expected output:

```
engine: camoufox browser: true
tab: <id> http://example.com/
result: Hi, friends!
```

The injected expression:

```javascript
document.title = "Hi, friends!"; document.title
```
