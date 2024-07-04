package model

type WhatsAppMessage struct {
	MessagingProduct string       `json:"messaging_product"`
	RecipientType    string       `json:"recipient_type"`
	To               string       `json:"to"`
	Type             string       `json:"type"`
	Text             MessageText  `json:"text"`
	Audio            MessageAudio `json:"audio"`
}

type MessageText struct {
	PreviewURL bool   `json:"preview_url"`
	Body       string `json:"body"`
}

type MessageAudio struct {
	To            string `json:"to"`
	Type          string `json:"type"`
	RecipientType string `json:"recipient_type"`
	Audio         Audio  `json:"audio"`
}

type Audio struct {
	ID string `json:"id"`
}
