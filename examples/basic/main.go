// Command basic connects to a running Camoufox instance (via the
// camofox-browser HTTP wrapper) and demonstrates tab lifecycle:
// open, navigate, snapshot, close.
//
// Run from the repository root:
//
//	go run ./examples/basic -addr http://localhost:9377
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	addr := flag.String("addr", "http://localhost:9377", "HTTP endpoint of the running camofox-browser")
	flag.Parse()

	client := &http.Client{Timeout: 30 * time.Second}
	userID := "demo"

	// 1. Health check
	{
		resp, err := client.Get(*addr + "/health")
		if err != nil {
			log.Fatalf("health: %v", err)
		}
		defer resp.Body.Close()
		var body map[string]any
		json.NewDecoder(resp.Body).Decode(&body)
		fmt.Printf("health: %v\n", body)
	}

	// 2. Open a tab
	var tabID string
	{
		payload := map[string]any{
			"userId":     userID,
			"sessionKey": "demo",
			"url":        "https://example.com",
		}
		tabID = postJSON(client, *addr+"/tabs/open", payload)
		fmt.Printf("tab opened: %s\n", tabID)
	}

	// 3. Navigate to another page
	{
		payload := map[string]any{
			"userId": userID,
			"url":    "https://httpbin.org/html",
		}
		postJSON(client, *addr+"/tabs/"+tabID+"/navigate", payload)
		fmt.Println("navigated to httpbin.org/html")
	}

	// 4. Get snapshot
	{
		url := fmt.Sprintf("%s/tabs/%s/snapshot?userId=%s", *addr, tabID, userID)
		resp, err := client.Get(url)
		if err != nil {
			log.Fatalf("snapshot: %v", err)
		}
		defer resp.Body.Close()
		var snap map[string]any
		json.NewDecoder(resp.Body).Decode(&snap)
		if s, ok := snap["snapshot"].(string); ok {
			fmt.Printf("snapshot (%d chars):\n%s\n", len(s), s)
		} else {
			fmt.Printf("snapshot: %v\n", snap)
		}
	}

	// 5. Close the tab
	{
		payload := map[string]any{"userId": userID}
		data, _ := json.Marshal(payload)
		req, _ := http.NewRequest("DELETE", *addr+"/tabs/"+tabID, bytes.NewReader(data))
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("close tab: %v", err)
		}
		resp.Body.Close()
		fmt.Println("tab closed")
	}

	// 6. Close session
	{
		req, _ := http.NewRequest("DELETE", *addr+"/sessions/"+userID, nil)
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("close session: %v", err)
		}
		resp.Body.Close()
		fmt.Println("session closed")
	}
}

func postJSON(client *http.Client, url string, payload map[string]any) string {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("marshal: %v", err)
	}
	resp, err := client.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		log.Fatalf("post %s: %v", url, err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		log.Fatalf("%s: %d %s", url, resp.StatusCode, body)
	}
	var result map[string]any
	json.Unmarshal(body, &result)
	if id, ok := result["tabId"].(string); ok {
		return id
	}
	if id, ok := result["targetId"].(string); ok {
		return id
	}
	log.Fatalf("no tabId in response: %s", body)
	return ""
}
