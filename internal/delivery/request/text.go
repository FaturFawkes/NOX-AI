package request

import "encoding/json"

type WhatsAppBusinessAccount struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

type Entry struct {
	ID      string   `json:"id"`
	Changes []Change `json:"changes"`
}

type Change struct {
	Value Value  `json:"value"`
	Field string `json:"field"`
}

type Value struct {
	MessagingProduct string    `json:"messaging_product"`
	Metadata         Metadata  `json:"metadata"`
	Contacts         []Contact `json:"contacts"`
	Messages         []Message `json:"messages"`
}

type Metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

type Contact struct {
	Profile Profile `json:"profile"`
	WaID    string  `json:"wa_id"`
}

type Profile struct {
	Name string `json:"name"`
}

type Message struct {
	Context     Context     `json:"context,omitempty"`
	From        string      `json:"from"`
	ID          string      `json:"id"`
	Timestamp   string      `json:"timestamp"`
	Text        Text        `json:"text"`
	Type        string      `json:"type"`
	Interactive Interactive `json:"interactive"`
}

type Text struct {
	Body string `json:"body"`
}

type Context struct {
	From string `json:"from"`
	ID   string `json:"id"`
}

type Interactive struct {
	Type      string    `json:"type"`
	ListReply ListReply `json:"list_reply"`
}

type ListReply struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type RequestData struct {
	Body json.RawMessage `json:"body"`
}
