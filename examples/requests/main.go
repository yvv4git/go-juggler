// Command requests demonstrates intercepting companion network requests
// (CSS, JS, images, trackers) loaded by a page.
//
// Run from the repository root:
//
//	go run ./examples/requests https://www.wikipedia.org
//	go run ./examples/requests -wait 20s https://rutube.ru
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/yvv4git/go-juggler"
)

func main() {
	addr := flag.String("addr", "http://localhost:9377", "camofox-browser endpoint")
	url := flag.String("url", "http://www.wikipedia.org", "page to load")
	session := flag.String("session", "requests-demo", "session key")
	wait := flag.Duration("wait", 15*time.Second, "how long to poll for late requests")

	flag.Parse()

	if flag.NArg() > 0 {
		*url = flag.Arg(0)
	}

	ctx := context.Background()
	c := juggler.NewClient(*addr)

	tab, err := c.OpenTab(ctx, *session, *url)
	if err != nil {
		log.Fatalf("OpenTab: %v", err)
	}

	fmt.Printf("opened %s\n\n", tab.URL)
	defer func() { _ = c.CloseTab(ctx, tab.TabID, *session) }()

	fmt.Printf("polling for %s ...\n\n", *wait)

	entries, err := c.PollNetworkRequests(ctx, tab.TabID, *session, *wait, 1*time.Second)
	if err != nil {
		log.Fatalf("PollNetworkRequests: %v", err)
	}

	fmt.Printf("total requests: %d\n\n", len(entries))

	for i, e := range entries {
		name := e.Name
		if len(name) > 100 {
			name = name[:97] + "..."
		}

		fmt.Printf("%2d. [%-12s] %s\n", i+1, strings.ToUpper(e.Type), name)
	}
}
