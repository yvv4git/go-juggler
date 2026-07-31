// Command screen demonstrates taking a screenshot of a page.
//
// Run from the repository root:
//
//	go run ./examples/screen https://www.wikipedia.org
//	go run ./examples/screen -out screenshot.png https://rutube.ru
//	go run ./examples/screen -full https://example.com
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/yvv4git/go-juggler"
)

func main() {
	addr := flag.String("addr", "http://localhost:9377", "camofox-browser endpoint")
	url := flag.String("url", "http://www.wikipedia.org", "page to load")
	session := flag.String("session", "screen-demo", "session key")
	out := flag.String("out", "screenshot.png", "output PNG file")
	fullPage := flag.Bool("full", false, "capture full page (not just viewport)")
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
	fmt.Printf("opened %s\n", tab.URL)
	defer c.CloseTab(ctx, tab.TabID, *session)

	png, err := c.Screenshot(ctx, tab.TabID, *session, *fullPage)
	if err != nil {
		log.Fatalf("Screenshot: %v", err)
	}

	if err := os.WriteFile(*out, png, 0644); err != nil {
		log.Fatalf("write file: %v", err)
	}
	fmt.Printf("saved %s (%d bytes)\n", *out, len(png))
}
