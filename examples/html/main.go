// Command html demonstrates fetching page HTML in two phases:
// first without waiting for dynamic content, then again after
// dynamic content has loaded.
//
// Run from the repository root:
//
//	go run ./examples/html https://www.wikipedia.org
//	go run ./examples/html -wait 5s https://rutube.ru
//	go run ./examples/html -addr http://localhost:9377 -wait 20s -session my-session https://example.com
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
	url := flag.String("url", "http://www.wikipedia.org", "page to load")
	session := flag.String("session", "html-demo", "session key")
	wait := flag.Duration("wait", 10*time.Second, "wait for dynamic content")

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
	defer func() { _ = c.CloseTab(ctx, tab.TabID, *session) }()

	// Phase 1: get HTML immediately (without waiting for dynamic content)
	early, err := c.Evaluate(ctx, tab.TabID, *session,
		`document.documentElement.outerHTML`)
	if err != nil {
		log.Fatalf("early evaluate: %v", err)
	}

	earlyHTML, _ := early.Result.(string)
	fmt.Printf("early HTML: %d bytes\n", len(earlyHTML))

	// Phase 2: wait for dynamic content, then re-fetch
	time.Sleep(*wait)

	late, err := c.Evaluate(ctx, tab.TabID, *session,
		`document.documentElement.outerHTML`)
	if err != nil {
		log.Fatalf("late evaluate: %v", err)
	}

	lateHTML, _ := late.Result.(string)
	fmt.Printf("late  HTML: %d bytes\n", len(lateHTML))

	// Use the later version if it changed
	if len(lateHTML) > len(earlyHTML) {
		fmt.Printf("dynamic content loaded: +%d bytes\n", len(lateHTML)-len(earlyHTML))
		earlyHTML = lateHTML
	}

	fmt.Println()
	fmt.Println(earlyHTML)
}
