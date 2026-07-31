# go-juggler

Go client for the Juggler
automation protocol. Automates Juggler-enabled browsers such as
Firefox/Camoufox from Go.

## Install

```sh
go get github.com/yvv4git/go-juggler
```

## Usage

```go
package main

import (
    "context"
    "log"

    "github.com/yvv4git/go-juggler"
)

func main() {
    ctx := context.Background()

    // Connect to a transport (pipe or WebSocket) or launch a browser.
    b, err := juggler.Launch(ctx, juggler.WithExecPath("/path/to/firefox"))
    if err != nil {
        log.Fatal(err)
    }
    defer b.Close(ctx)

    // tabs, navigation, and the rest of the API are on their way.
}
```

## Layout

The module is a library split into three layers; the root package
re-exports the public browser API so most callers import a single package.

```text
go-juggler/
├── juggler.go          # root package: public API facade
├── transport/          # byte channels: pipe (fd 3/fd 4) and WebSocket
├── protocol/           # Juggler message types, JSON encoding/decoding
└── browser/            # high-level browser and tab control
```

Dependencies flow one way, from high level to low level:

`browser/` → `transport/`, and `protocol/` stays independent so transport
moves raw frames and protocol gives them meaning.

## Links

- [Juggler Protocol Architecture (daijro/camoufox)](https://deepwiki.com/daijro/camoufox/6.1-juggler-system)

## License

MIT, see [LICENSE](LICENSE). See [NOTICE](NOTICE) for attribution.
