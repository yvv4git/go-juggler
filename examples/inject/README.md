# examples/inject

Demonstrates injecting JavaScript into a page via `POST /tabs/:tabId/evaluate`.

**Requires camofox-browser >= 1.4.0**

```bash
CAMOFOX_SESSION=demo go run .
```

Expected output:

```
engine: camoufox browser: true
tab: <id> https://httpbin.org/html
result: done
```

The injected script is:

```javascript
console.log("Hi, friends!");
"done"
```
