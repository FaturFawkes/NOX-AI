package model

type ImageMessage struct {
	MessagingProduct string  `json:"messaging_product"`
	RecipientType    string  `json:"recipient_type"`
	To               string  `json:"to"`
	Type             string  `json:"type"`
	Image            Image   `json:"image"`
}