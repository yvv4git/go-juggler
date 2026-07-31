# examples/html

Demonstrates fetching page HTML in two phases:
first without waiting for dynamic content, then again after dynamic content has loaded.

```bash
go run . -wait 5s https://www.wikipedia.org
```

## Flags

| Flag       | Default                     | Description                     |
|------------|-----------------------------|---------------------------------|
| `-addr`    | `http://localhost:9377`     | camofox-browser endpoint        |
| `-url`     | `http://www.wikipedia.org`  | page to load (or positional)    |
| `-session` | `html-demo`                 | session key                     |
| `-wait`    | `10s`                       | wait for dynamic content        |

## Sample output

```text
early HTML: 121000 bytes
late  HTML: 121051 bytes
dynamic content loaded: +51 bytes

<html lang="en" class="js-enabled"><head>
...
```
