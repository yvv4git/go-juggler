// Command tab demonstrates every tab operation available in go-juggler:
// open, navigate, snapshot, click, type, press, scroll, back, forward,
// refresh, links, screenshot, stats, and close.
//
// Run from the repository root:
//
//	go run ./examples/tab -addr http://localhost:9377
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
	flag.Parse()

	ctx := context.Background()
	c := juggler.NewClient(*addr)
	session := "tab-demo"

	// 1. OpenTab - open a blank tab, then navigate separately
	tab, err := c.OpenTab(ctx, session, "https://example.com")
	if err != nil {
		log.Fatalf("OpenTab: %v", err)
	}
	fmt.Printf("1. OpenTab    -> %s (%s)\n", tab.TabID, tab.URL)
	defer func() {
		_ = c.CloseTab(ctx, tab.TabID, session)
		_ = c.CloseSession(ctx, session)
		fmt.Println("\n14. CloseTab    ")
		fmt.Println("15. CloseSession ")
	}()

	// 2. Navigate - load a real page with interactive elements
	if err := c.Navigate(ctx, tab.TabID, session, "https://www.wikipedia.org"); err != nil {
		log.Fatalf("Navigate: %v", err)
	}
	fmt.Println("2. Navigate   -> wikipedia.org")

	// 3. Snapshot - get the ARIA tree to find refs
	snap, err := c.Snapshot(ctx, tab.TabID, session)
	if err != nil {
		log.Fatalf("Snapshot: %v", err)
	}
	fmt.Printf("3. Snapshot   -> %d chars, %d refs\n", len(snap.Snapshot), snap.RefsCount)

	// 4. Links - list all links on the page
	links, err := c.Links(ctx, tab.TabID, session, 5, 0)
	if err != nil {
		log.Fatalf("Links: %v", err)
	}
	fmt.Printf("4. Links      -> %d total, first 5:\n", links.Pagination.Total)
	for _, l := range links.Links {
		fmt.Printf("   - %s (%s)\n", l.Text, l.URL)
	}

	// 5. Stats - tab state
	stats, err := c.Stats(ctx, tab.TabID, session)
	if err != nil {
		log.Fatalf("Stats: %v", err)
	}
	fmt.Printf("5. Stats      -> url=%s visited=%d refs=%d\n", stats.URL, len(stats.VisitedURLs), stats.RefsCount)

	// 6. Type - search for something on Wikipedia
	//    The search box is typically ref e3 on wikipedia.org
	if err := c.Type(ctx, tab.TabID, session, "e3", "", "Go programming language"); err != nil {
		fmt.Printf("6. Type       -> skipped (ref e3 not found on this page): %v\n", err)
	} else {
		fmt.Println("6. Type       -> \"Go programming language\" into search box")
	}

	// 7. Press - press Enter to submit the search
	if err := c.Press(ctx, tab.TabID, session, "Enter"); err != nil {
		fmt.Printf("7. Press      -> skipped: %v\n", err)
	} else {
		fmt.Println("7. Press      -> Enter")
	}

	// 8. Scroll - scroll down 500px
	if err := c.Scroll(ctx, tab.TabID, session, "down", 500); err != nil {
		fmt.Printf("8. Scroll     -> %v\n", err)
	} else {
		fmt.Println("8. Scroll     -> down 500px")
	}

	// 9. Back - go back to the previous page
	if err := c.Back(ctx, tab.TabID, session); err != nil {
		fmt.Printf("9. Back       -> %v\n", err)
	} else {
		fmt.Println("9. Back       -> returned to previous page")
	}

	// 10. Forward - go forward again
	if err := c.Forward(ctx, tab.TabID, session); err != nil {
		fmt.Printf("10. Forward   -> %v\n", err)
	} else {
		fmt.Println("10. Forward   -> moved forward")
	}

	// 11. Refresh - reload the page
	if err := c.Refresh(ctx, tab.TabID, session); err != nil {
		fmt.Printf("11. Refresh   -> %v\n", err)
	} else {
		fmt.Println("11. Refresh   -> reloaded")
	}

	// 12. Click - click a link by ref (e.g. "e1" from snapshot)
	if err := c.Click(ctx, tab.TabID, session, "e1", ""); err != nil {
		fmt.Printf("12. Click     -> skipped (ref e1 may not exist): %v\n", err)
	} else {
		fmt.Println("12. Click     -> clicked e1")
	}

	// 13. Screenshot - save a PNG
	png, err := c.Screenshot(ctx, tab.TabID, session, false)
	if err != nil {
		fmt.Printf("13. Screenshot -> %v\n", err)
	} else {
		path := "screenshot.png"
		_ = os.WriteFile(path, png, 0644)
		fmt.Printf("13. Screenshot -> saved %s (%d bytes)\n", path, len(png))
	}
}
