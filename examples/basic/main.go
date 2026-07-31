// Command basic connects to a running Camoufox instance (via the
// camofox-browser HTTP wrapper) and demonstrates the full tab lifecycle
// using the go-juggler library.
//
// Run from the repository root:
//
//	go run ./examples/basic -addr http://localhost:9377
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/yvv4git/go-juggler"
)

func main() {
	addr := flag.String("addr", "http://localhost:9377", "HTTP endpoint of the running camofox-browser")
	flag.Parse()

	ctx := context.Background()

	c := juggler.NewClient(*addr)

	// 1. Health check
	health, err := c.Health(ctx)
	if err != nil {
		log.Fatalf("health: %v", err)
	}
	fmt.Printf("health: ok=%v engine=%s sessions=%d\n", health.OK, health.Engine, health.Sessions)

	// 2. Open a tab
	tab, err := c.OpenTab(ctx, "demo", "https://example.com")
	if err != nil {
		log.Fatalf("open tab: %v", err)
	}
	fmt.Printf("tab opened: %s\n", tab.TabID)

	// 3. Navigate
	if err := c.Navigate(ctx, tab.TabID, "demo", "https://rutube.ru"); err != nil {
		log.Fatalf("navigate: %v", err)
	}
	fmt.Println("navigated to rutube.ru")

	// 4. Snapshot
	snap, err := c.Snapshot(ctx, tab.TabID, "demo")
	if err != nil {
		log.Fatalf("snapshot: %v", err)
	}
	fmt.Printf("snapshot (%d chars):\n%s\n", len(snap.Snapshot), snap.Snapshot)

	// 5. Close tab
	if err := c.CloseTab(ctx, tab.TabID, "demo"); err != nil {
		log.Fatalf("close tab: %v", err)
	}
	fmt.Println("tab closed")

	// 6. Close session
	if err := c.CloseSession(ctx, "demo"); err != nil {
		log.Fatalf("close session: %v", err)
	}
	fmt.Println("session closed")
}
