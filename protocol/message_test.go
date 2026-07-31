package protocol

import (
	"encoding/json"
	"testing"
)

func TestRequestRoundTrip(t *testing.T) {
	req, err := NewRequest(1, BrowserNewPage, map[string]any{"url": "https://example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if !req.IsRequest() {
		t.Fatalf("expected request, got %+v", req)
	}

	data, err := req.Encode()
	if err != nil {
		t.Fatal(err)
	}

	got, err := Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != 1 || got.Method != BrowserNewPage {
		t.Fatalf("unexpected message: %+v", got)
	}
	var params map[string]string
	if err := json.Unmarshal(got.Params, &params); err != nil {
		t.Fatalf("decode params: %v", err)
	}
	if params["url"] != "https://example.com" {
		t.Fatalf("unexpected params: %v", params)
	}
}

func TestEventRoundTrip(t *testing.T) {
	data, err := json.Marshal(&Message{Method: "Page.crashed", Params: json.RawMessage(`{}`)})
	if err != nil {
		t.Fatal(err)
	}
	msg, err := Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if !msg.IsEvent() {
		t.Fatalf("expected event, got %+v", msg)
	}
	if msg.IsRequest() || msg.IsResponse() {
		t.Fatalf("event misclassified: %+v", msg)
	}
}

func TestResponseClassification(t *testing.T) {
	data, err := json.Marshal(&Message{ID: 7, Result: json.RawMessage(`{"ok":true}`)})
	if err != nil {
		t.Fatal(err)
	}
	msg, err := Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if !msg.IsResponse() {
		t.Fatalf("expected response, got %+v", msg)
	}

	var result struct{ Ok bool }
	if err := json.Unmarshal(msg.Result, &result); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	if !result.Ok {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestErrorString(t *testing.T) {
	err := &Error{Code: 123, Message: "boom"}
	if got, want := err.Error(), "juggler: boom (123)"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
