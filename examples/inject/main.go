package main

import (
	"context"
	"fmt"
	"log"
	"os"

	juggler "github.com/yvv4git/go-juggler"
)

func main() {
	addr := os.Getenv("CAMOFOX_ADDR")
	if addr == "" {
		addr = "http://localhost:9377"
	}
	session := os.Getenv("CAMOFOX_SESSION")
	if session == "" {
		session = "demo"
	}

	client := juggler.NewClient(addr)

	health, err := client.Health(context.Background())
	if err != nil {
		log.Fatalf("health: %v", err)
	}
	fmt.Println("engine:", health.Engine, "browser:", health.BrowserConnected)

	tab, err := client.OpenTab(context.Background(), session, "https://httpbin.org/html")
	if err != nil {
		log.Fatalf("open tab: %v", err)
	}
	fmt.Println("tab:", tab.TabID, tab.URL)
	defer client.CloseTab(context.Background(), tab.TabID, session)

	res, err := client.Evaluate(context.Background(), tab.TabID, session,
		`console.log("Hi, friends!"); "done"`)
	if err != nil {
		log.Fatalf("evaluate: %v\nHint: camofox-browser must be >= 1.4.0", err)
	}
	fmt.Println("result:", res.Result)
}
