package protocol

import (
	"encoding/json"
	"strconv"
)

// Message is a single frame exchanged between the client and the browser.
//
// A request carries an ID and a Method. A response carries the ID of the
// request it answers and either a Result or an Error. An event carries a
// Method but no ID.
type Message struct {
	ID     int64           `json:"id,omitempty"`
	Method string          `json:"method,omitempty"`
	Params json.RawMessage `json:"params,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *Error          `json:"error,omitempty"`
}

// Error describes a protocol-level failure reported by the browser.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return "juggler: " + e.Message + " (" + strconv.Itoa(e.Code) + ")"
}

// IsRequest reports whether m is a client-to-browser request.
func (m *Message) IsRequest() bool { return m.Method != "" && m.ID != 0 }

// IsResponse reports whether m answers a previous request.
func (m *Message) IsResponse() bool { return m.Method == "" && m.ID != 0 }

// IsEvent reports whether m is a browser-initiated notification.
func (m *Message) IsEvent() bool { return m.Method != "" && m.ID == 0 }

// NewRequest builds a request message for the given method, encoding params
// as JSON.
func NewRequest(id int64, method string, params any) (*Message, error) {
	raw, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	return &Message{ID: id, Method: method, Params: raw}, nil
}

// Decode parses a Message from wire-encoded JSON.
func Decode(data []byte) (*Message, error) {
	var m Message
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// Encode serialises the message to wire-encoded JSON.
func (m *Message) Encode() ([]byte, error) {
	return json.Marshal(m)
}
