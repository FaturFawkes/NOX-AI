package model

type InteractiveMessage struct {
	MessagingProduct string          `json:"messaging_product"`
	RecipientType    string          `json:"recipient_type"`
	To               string          `json:"to"`
	Type             string          `json:"type"`
	Interactive      InteractiveData `json:"interactive"`
}

type InteractiveData struct {
	Type   string            `json:"type"`
	Body   InteractiveText   `json:"body"`
	Action InteractiveAction `json:"action"`
}

type InteractiveText struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}

type InteractiveAction struct {
	Button   string               `json:"button"`
	Sections []InteractiveSection `json:"sections"`
}

type InteractiveSection struct {
	Title string           `json:"title"`
	Rows  []InteractiveRow `json:"rows"`
}

type InteractiveRow struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
