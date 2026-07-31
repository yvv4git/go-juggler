// Package protocol implements the Juggler wire protocol: message types,
// commands and their JSON encoding, as exchanged between a client and a
// Juggler-enabled browser (Firefox/Camoufox).
//
// The wire format follows a CDP-style message layout. Requests carry an id
// and a method, responses match the request id and carry either a result or
// an error, and events are notifications without an id:
//
//	{ "id": 1, "method": "Browser.newPage", "params": { ... } }
//	{ "id": 1, "result": { ... } }
//	{ "method": "Page.crashed", "params": { ... } }
package protocol
