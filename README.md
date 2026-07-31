# go-juggler

Go client for the [Juggler](https://deepwiki.com/daijro/camoufox/6.1-juggler-system)
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

    b, err := juggler.Launch(ctx, juggler.WithExecPath("/path/to/firefox"))
    if err != nil {
        log.Fatal(err)
    }
    defer b.Close(ctx)
}
```

## Layout

```text
go-juggler/
├── juggler.go          # root package: public API facade
├── transport/          # byte channels: pipe (fd 3/fd 4) and WebSocket
├── protocol/           # Juggler message types, JSON encoding/decoding
├── browser/            # high-level browser and tab control
└── examples/           # runnable example programs
    └── basic/          # go run ./examples/basic -exec /path/to/firefox
```

Dependencies flow one way, from high level to low level:

`browser/` -> `transport/`, and `protocol/` stays independent so transport
moves raw frames and protocol gives them meaning.

## Links

- [Juggler Protocol Architecture (daijro/camoufox)](https://deepwiki.com/daijro/camoufox/6.1-juggler-system)

## License

MIT, see [LICENSE](LICENSE). See [NOTICE](NOTICE) for attribution.
