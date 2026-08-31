// Command headers demonstrates how to use custom HTTP headers
// with the go-juggler library.
//
// Run from the repository root:
//
//	go run ./examples/headers -addr http://localhost:9377
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

	// 2. Open a tab with custom headers
	headers := []juggler.Header{
		{Name: "Authorization", Value: "Bearer my-secret-token"},
		{Name: "X-Custom-Header", Value: "custom-value"},
		{Name: "Accept-Language", Value: "en-US,en;q=0.9"},
	}

	tab, err := c.OpenTab(ctx, "demo", "https://httpbin.org/headers", headers...)
	if err != nil {
		log.Fatalf("open tab: %v", err)
	}

	fmt.Printf("tab opened: %s\n", tab.TabID)

	// 3. Snapshot to see the page content
	snap, err := c.Snapshot(ctx, tab.TabID, "demo")
	if err != nil {
		log.Fatalf("snapshot: %v", err)
	}

	fmt.Printf("snapshot (%d chars):\n%s\n", len(snap.Snapshot), snap.Snapshot)

	// 4. Navigate to another URL with different headers
	navigateHeaders := []juggler.Header{
		{Name: "Authorization", Value: "Bearer another-token"},
		{Name: "X-Request-ID", Value: "req-12345"},
	}

	if err := c.Navigate(ctx, tab.TabID, "demo", "https://httpbin.org/headers", navigateHeaders...); err != nil {
		log.Fatalf("navigate: %v", err)
	}

	fmt.Println("navigated with custom headers")

	// 5. Snapshot the new page
	snap2, err := c.Snapshot(ctx, tab.TabID, "demo")
	if err != nil {
		log.Fatalf("snapshot: %v", err)
	}

	fmt.Printf("snapshot (%d chars):\n%s\n", len(snap2.Snapshot), snap2.Snapshot)

	// 6. Close tab
	if err := c.CloseTab(ctx, tab.TabID, "demo"); err != nil {
		log.Fatalf("close tab: %v", err)
	}

	fmt.Println("tab closed")

	// 7. Close session
	if err := c.CloseSession(ctx, "demo"); err != nil {
		log.Fatalf("close session: %v", err)
	}

	fmt.Println("session closed")
}
