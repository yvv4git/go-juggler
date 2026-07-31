// Command tabs demonstrates working with multiple tabs: open several,
// list them, and close them all.
//
// Run from the repository root:
//
//	go run ./examples/tabs
//	go run ./examples/tabs -count 5
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/yvv4git/go-juggler"
)

func main() {
	addr := flag.String("addr", "http://localhost:9377", "camofox-browser endpoint")
	session := flag.String("session", "tabs-demo", "session key")
	count := flag.Int("count", 3, "number of tabs to open")

	flag.Parse()

	ctx := context.Background()
	c := juggler.NewClient(*addr)

	urls := []string{
		"http://example.com",
		"http://example.org",
		"http://example.net",
		"http://example.edu",
		"http://example.gov",
	}

	for i := 0; i < *count; i++ {
		u := urls[i%len(urls)]

		tab, err := c.OpenTab(ctx, *session, u)
		if err != nil {
			log.Fatalf("OpenTab %d: %v", i+1, err)
		}

		fmt.Printf("opened %d: %s\n", i+1, tab.URL)
	}

	time.Sleep(1 * time.Second)

	r, err := c.ListTabs(ctx, *session)
	if err != nil {
		log.Fatalf("ListTabs: %v", err)
	}

	fmt.Printf("\ntotal tabs: %d\n\n", len(r.Tabs))

	for i, t := range r.Tabs {
		fmt.Printf("  %d. [%s] %s\n", i+1, t.Title, t.URL)
	}

	fmt.Println()

	for _, t := range r.Tabs {
		if err := c.CloseTab(ctx, t.TabID, *session); err != nil {
			fmt.Printf("close %s: %v\n", t.TabID[:8], err)
		} else {
			fmt.Printf("closed %s\n", t.TabID[:8])
		}
	}

	_ = c.CloseSession(ctx, *session)

	fmt.Println("\nsession closed")
}
