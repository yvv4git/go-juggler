// Command tab demonstrates every tab operation available in go-juggler:
// open, navigate, snapshot, click, type, press, scroll, back, forward,
// refresh, links, screenshot, stats, evaluate, and close.
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
	"time"

	"github.com/yvv4git/go-juggler"
)

func main() {
	addr := flag.String("addr", "http://localhost:9377", "camofox-browser endpoint")

	flag.Parse()

	ctx := context.Background()
	c := juggler.NewClient(*addr)
	session := "tab-demo"

	// 1. OpenTab
	tab, err := c.OpenTab(ctx, session, "http://example.com")
	if err != nil {
		log.Fatalf("OpenTab: %v", err)
	}

	fmt.Printf("1. OpenTab    -> %s (%s)\n", tab.TabID, tab.URL)
	defer func() {
		_ = c.CloseTab(ctx, tab.TabID, session)
		_ = c.CloseSession(ctx, session)

		fmt.Println("\n15. CloseTab    ")
		fmt.Println("16. CloseSession")
	}()

	// 2. Snapshot
	snap, err := c.Snapshot(ctx, tab.TabID, session)
	if err != nil {
		log.Fatalf("Snapshot: %v", err)
	}

	fmt.Printf("2. Snapshot   -> %d chars, %d refs\n", len(snap.Snapshot), snap.RefsCount)

	// 3. Links
	links, err := c.Links(ctx, tab.TabID, session, 5, 0)
	if err != nil {
		log.Fatalf("Links: %v", err)
	}

	fmt.Printf("3. Links      -> %d total, first 5:\n", links.Pagination.Total)

	for _, l := range links.Links {
		fmt.Printf("   - %s (%s)\n", l.Text, l.URL)
	}

	// 4. Stats
	stats, err := c.Stats(ctx, tab.TabID, session)
	if err != nil {
		log.Fatalf("Stats: %v", err)
	}

	fmt.Printf("4. Stats      -> url=%s visited=%d refs=%d\n", stats.URL, len(stats.VisitedURLs), stats.RefsCount)

	// 5. Evaluate - inject a clickable button, then Click it
	if _, err := c.Evaluate(ctx, tab.TabID, session,
		`document.body.innerHTML = "<button id='btn' onclick='document.title=\"clicked\"'>Click me</button>"; "ok"`); err != nil {
		fmt.Printf("5. Evaluate   -> skipped: %v\n", err)
	} else {
		fmt.Println("5. Evaluate   -> injected button")
	}

	snapBtn, err := c.Snapshot(ctx, tab.TabID, session)
	if err != nil {
		fmt.Printf("5. Snapshot   -> %v\n", err)
	} else {
		fmt.Printf("5. Snapshot   -> %d refs\n", snapBtn.RefsCount)
	}

	if err := c.Click(ctx, tab.TabID, session, "e1", ""); err != nil {
		fmt.Printf("5. Click      -> skipped: %v\n", err)
	} else {
		fmt.Println("5. Click      -> clicked e1 (button)")
	}

	// 6. Navigate - page B for Back/Forward
	if err := c.Navigate(ctx, tab.TabID, session, "http://example.org"); err != nil {
		log.Fatalf("Navigate: %v", err)
	}

	fmt.Println("6. Navigate   -> example.org")

	// 7. Back
	if err := c.Back(ctx, tab.TabID, session); err != nil {
		fmt.Printf("7. Back       -> skipped: %v\n", err)
	} else {
		fmt.Println("7. Back       -> returned to page A")
	}

	// 8. Forward
	if err := c.Forward(ctx, tab.TabID, session); err != nil {
		fmt.Printf("8. Forward    -> skipped: %v\n", err)
	} else {
		fmt.Println("8. Forward    -> moved to page B")
	}

	// 9. Evaluate - inject a form for Type test
	if _, err := c.Evaluate(ctx, tab.TabID, session,
		`document.body.innerHTML = "<form><input id='q' type='text' placeholder='Search...'><button type='submit'>OK</button></form>"; "injected"`); err != nil {
		fmt.Printf("9. Evaluate   -> skipped: %v\n", err)
	} else {
		fmt.Println("9. Evaluate   -> injected form")
	}

	// 10. Snapshot
	snap2, err := c.Snapshot(ctx, tab.TabID, session)
	if err != nil {
		log.Fatalf("Snapshot: %v", err)
	}

	fmt.Printf("10. Snapshot  -> %d chars, %d refs\n", len(snap2.Snapshot), snap2.RefsCount)

	// 11. Type
	if err := c.Type(ctx, tab.TabID, session, "e1", "", "Go lang"); err != nil {
		fmt.Printf("11. Type      -> skipped: %v\n", err)
	} else {
		fmt.Println("11. Type      -> \"Go lang\" into search box")
	}

	// 12. Press
	if err := c.Press(ctx, tab.TabID, session, "Enter"); err != nil {
		fmt.Printf("12. Press     -> skipped: %v\n", err)
	} else {
		fmt.Println("12. Press     -> Enter")
	}

	time.Sleep(1 * time.Second)

	// 13. Scroll
	if err := c.Scroll(ctx, tab.TabID, session, "down", 500); err != nil {
		fmt.Printf("13. Scroll    -> %v\n", err)
	} else {
		fmt.Println("13. Scroll    -> down 500px")
	}

	// 14. Refresh
	if err := c.Refresh(ctx, tab.TabID, session); err != nil {
		fmt.Printf("14. Refresh   -> %v\n", err)
	} else {
		fmt.Println("14. Refresh   -> reloaded")
	}

	// 15. Screenshot
	png, err := c.Screenshot(ctx, tab.TabID, session, false)
	if err != nil {
		fmt.Printf("15. Screenshot -> %v\n", err)
	} else {
		path := "screenshot.png"
		_ = os.WriteFile(path, png, 0o644)
		fmt.Printf("15. Screenshot -> saved %s (%d bytes)\n", path, len(png))
	}
}
