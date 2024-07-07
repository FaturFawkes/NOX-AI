package model

type WhatsAppMessage struct {
	MessagingProduct string         `json:"messaging_product"`
	RecipientType    string         `json:"recipient_type"`
	To               string         `json:"to"`
	Type             string         `json:"type"`
	Text             MessageText    `json:"text,omitempty"`
	Audio            MessageAudio   `json:"audio,omitempty"`
	Context          ContextMessage `json:"context,omitempty"`
}

type MessageText struct {
	PreviewURL bool   `json:"preview_url"`
	Body       string `json:"body"`
}

type MessageAudio struct {
	ID string `json:"id"`
}

type Audio struct {
	ID string `json:"id"`
}

type ContextMessage struct {
	MessageID string `json:"message_id"`
}
