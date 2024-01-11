package model

type Image struct {
	Link string `json:"link"`
}

type Parameter struct {
	Type  string `json:"type"`
	Image Image  `json:"image"`
}

type Component struct {
	Type        string      `json:"type"`
	Parameters  []Parameter `json:"parameters"`
}

type Language struct {
	Code string `json:"code"`
}

type Template struct {
	Name       string      `json:"name"`
	Language   Language    `json:"language"`
	Components []Component `json:"components"`
}

type MessageTemplate struct {
	MessagingProduct string   `json:"messaging_product"`
	RecipientType    string   `json:"recipient_type"`
	To               string   `json:"to"`
	Type             string   `json:"type"`
	Template         Template `json:"template"`
}
