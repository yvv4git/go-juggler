# examples/screen

Demonstrates taking a screenshot of a page and saving it as PNG.

```bash
go run . http://example.com
go run . -out screenshot.png -full http://www.wikipedia.org
```

## Flags

| Flag       | Default                     | Description                     |
|------------|-----------------------------|---------------------------------|
| `-addr`    | `http://localhost:9377`     | camofox-browser endpoint        |
| `-url`     | `http://www.wikipedia.org`  | page to load (or positional)    |
| `-session` | `screen-demo`               | session key                     |
| `-out`     | `screenshot.png`            | output PNG file                 |
| `-full`    | `false`                     | capture full page (not viewport)|

## Sample output

```text
opened http://example.com/
saved screenshot.png (36846 bytes)
```
