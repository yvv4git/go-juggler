# examples/tabs

Demonstrates working with multiple tabs: open several, list them, and close them.

```bash
go run .
go run . -count 5
```

## Flags

| Flag       | Default                     | Description                  |
|------------|-----------------------------|------------------------------|
| `-addr`    | `http://localhost:9377`     | camofox-browser endpoint     |
| `-session` | `tabs-demo`                 | session key                  |
| `-count`   | `3`                         | number of tabs to open       |

## Sample output

```text
opened 1: http://example.com/
opened 2: http://example.org/
opened 3: http://example.net/

total tabs: 3

  1. [Example Domain] http://example.com/
  2. [Example Domain] http://example.org/
  3. [Example Domain] http://example.net/

closed acd74c0c
closed 8c0951ac
closed 57adc6dd

session closed
```
